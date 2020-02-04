package modules

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/event"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util/threads"
	"golang.org/x/net/websocket"
)

func WebSocketSetIpPort(ipPort string) options {
	return func(mod IModule) {
		mod.(*WebSocketModule).ipPort = ipPort
	}
}

//超时时间（秒）
//例：超时时间为10秒时，就传入10
func WebSocketSetTimeout(timeout time.Duration) options {
	return func(mod IModule) {
		mod.(*WebSocketModule).timeout = timeout * time.Second
	}
}

//连接成功后回调，可以用来获取一些连接的信息，比如IP
func WebScoketSetOnlineFun(fun func(conn *event.WebSocketModel)) options {
	return func(mod IModule) {
		mod.(*WebSocketModule).webSocketOnlineFun = fun
	}
}

type WebSocketModule struct {
	ipPort             string                           //HTTP监听的地址
	timeout            time.Duration                    //超时时间
	RouteHandle        messages.IMessageHandle          //消息路由
	webSocketOnlineFun func(conn *event.WebSocketModel) //连接成功后回调，可以用来获取一些连接的信息，比如IP
	getnum             int64                            //收到的总消息数
	runing             int64                            //当前在处理的消息数
	connlen            int64                            //连接数
	httpServer         *http.Server                     //HTTP请求的对象
	thgo               *threads.ThreadGo                //协程管理器
}

func NewWebSocketModule(opts ...options) *WebSocketModule {
	result := &WebSocketModule{
		ipPort:             ":8081",
		timeout:            60 * time.Second,
		getnum:             0,
		runing:             0,
		connlen:            0,
		thgo:               threads.NewThreadGo(),
		RouteHandle:        messages.JsonMessageHandleNew(),
		webSocketOnlineFun: nil,
	}

	for _, opt := range opts {
		opt(result)
	}
	return result
}

//Init 初始化
func (mod *WebSocketModule) Init() {
	mod.httpServer = &http.Server{
		Addr:         mod.ipPort,
		WriteTimeout: mod.timeout * 2,
	}
	//还可以加别的参数，已后再加，有需要再加
	mux := http.NewServeMux()
	//这个是主要的逻辑
	mux.Handle("/", websocket.Handler(mod.Handle))
	//你也可以在外面继续扩展

	mod.httpServer.Handler = mux
}

//Start 启动
func (mod *WebSocketModule) Start() {
	//启动的协程
	mod.thgo.Go(func(ctx context.Context) {
		Logger.PStatus("websocket Module Start!")
		err := mod.httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				Logger.PStatus("websocket run Server closed under requeset!!")
				// log.Print("Server closed under requeset!!")
			} else {
				Logger.PFatal("Server closed unexpecteed:" + err.Error())
				// log.Fatal("Server closed unexpecteed!!")
			}
		}
	})
}

//Stop 停止
func (mod *WebSocketModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		Logger.PError(err, "Close websocket Module:")
	}
	mod.thgo.CloseWait()
	Logger.PStatus("websocket Module Stop.")
}

//PrintStatus 打印状态
func (mod *WebSocketModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\twebsocket Module\t:%d/%d/%d\t(connum/getmsg/runing)",
		atomic.LoadInt64(&mod.connlen),
		atomic.LoadInt64(&mod.getnum),
		atomic.LoadInt64(&mod.runing))
}

//Handle http发来的所有请求都会到这个方法来
func (mod *WebSocketModule) Handle(conn *websocket.Conn) {

	//标注子连接是不是都停下来
	mod.thgo.Wg.Add(1)
	defer mod.thgo.Wg.Done()
	defer conn.Close()

	//发给下面的连接对象，可以自定义一些信息和回调
	wsconn := new(event.WebSocketModel)
	wsconn.Conn = conn
	wsconn.KeyID = -1
	if mod.webSocketOnlineFun != nil {
		mod.webSocketOnlineFun(wsconn)
	}
	atomic.AddInt64(&mod.connlen, 1)
	//发消息来说明这个用户掉线了
	defer func() {
		atomic.AddInt64(&mod.connlen, -1)
		Logger.PDebugKey("websocket client closeing:%+v .", wsconn.KeyID, wsconn.ConInfo)
		//用来处理发生连接关闭的时候，要处理的事
		if wsconn.CloseFun != nil {
			wsconn.CloseFun(wsconn)
		}
		Logger.PDebugKey("websocket client close%+v .", wsconn.KeyID, wsconn.ConInfo)
	}()
	Logger.PDebugKey("websocket client open%+v .", wsconn.KeyID, wsconn.ConInfo)
	runchan := make(chan bool, 8) //用来处理超时
	mod.thgo.Go(
		func(ctx context.Context) {
			timeout := time.NewTimer(mod.timeout)
			defer timeout.Stop()
			defer conn.Close()
			for {
				select {
				case <-ctx.Done():
					return
				case <-timeout.C:
					return
				case ok := <-runchan:
					if ok {
						timeout.Reset(mod.timeout)
					} else {
						return
					}
				}
			}
			//超时关连接
		})
	mod.thgo.Try(
		func(ctx context.Context) {
		listen:
			for {
				buff, err := ioutil.ReadAll(conn)
				if err != nil {
					if err == io.EOF {
						runchan <- false
					}
					break listen
				}

				msg, err := mod.RouteHandle.Unmarshal(buff)
				if err != nil {
					Logger.PInfo("RouteHandle Unmarshal Error:%s", err.Error())
					return
				}
				modmsg, ok := msg.(messages.IWebSocketMessageHandle)
				if !ok {
					Logger.PInfo("Not is Web socket Msg:%+v", msg)
					return
				} else {
					Logger.PInfo("Web socket Get Msg:%+v", msg)
				}

				runchan <- true
				atomic.AddInt64(&mod.getnum, 1)
				mod.thgo.Try(func(ctx context.Context) {
					atomic.AddInt64(&mod.runing, 1)
					modmsg.WebSocketDirectCall(wsconn)
				}, nil, func() {
					atomic.AddInt64(&mod.runing, -1)
				})

			}
		},
		nil,
		nil,
	)

}

//GetPlayerNum用户连接数量
func (mod *WebSocketModule) GetPlayerNum() int64 {
	return atomic.LoadInt64(&mod.connlen)
}
