package module

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/event"
	"encoding/json"
	"github.com/buguang01/util/threads"

	"golang.org/x/net/websocket"
)

//如果一个服务器，又要接受客户端的，又要接受其他服务器的消息转发
//建议开二个这样的服务器，然后用一个BridgeModule,做消息的转发
//

//WebSocketConfig websocket的配置
type WebSocketConfig struct {
	Addr      string //websocket 监听地址
	Timeout   int    //超时时间 （秒）
	MsgMaxLen int    //一个消息最大长度B

}

//WebSocketModule socket监听模块
type WebSocketModule struct {
	Addr       string            //HTTP监听的地址
	httpServer *http.Server      //HTTP请求的对象
	cg         *WebSocketConfig  //从配置表中读进来的数据
	wg         sync.WaitGroup    //用来确定是不是关闭了
	threadgo   *threads.ThreadGo //子协程管理
	getnum     int64             //收到的总消息数
	sendnum    int64             //当前在处理的消息数

	wsmap     map[*websocket.Conn]bool //所有的连接
	wsmaplock sync.Mutex               //上面那个对象的锁

	RouteFun           func(code int) event.WebSocketCall //用来生成事件处理器的工厂
	WebSocketOnlineFun func(conn *websocket.Conn) string  //连接成功后
	// wslist         map[uint]*WsConnModel    //websocket 列表
	// wsclosefunlist map[*websocket.Conn]func() //如果socket关闭的时候调用
	// wslock         sync.RWMutex               //对上面那个列表的锁

	// ConnCloseFun event.WebSocketCall                  //连接断开的时候
	// TimeoutFun func(map[string]interface{}) []byte  //超时时的回调方法
}

//NewWSModule 生成一个新的websocket的对象
func NewWSModule(configmd *WebSocketConfig) *WebSocketModule {
	result := &WebSocketModule{
		cg:   configmd,
		Addr: configmd.Addr,
	}
	result.wsmap = make(map[*websocket.Conn]bool)
	result.threadgo = threads.NewThreadGo()
	// result.wslist = make(map[uint]*WsConnModel)
	// result.wsclosefunlist = make(map[*websocket.Conn]func())
	return result
}

//Init IModule接口的实现
func (mod *WebSocketModule) Init() {

	mod.httpServer = &http.Server{
		Addr:         mod.cg.Addr,
		WriteTimeout: time.Duration(mod.cg.Timeout) * time.Second,
	}
	//还可以加别的参数，已后再加，有需要再加
	mux := http.NewServeMux()
	//这个是主要的逻辑
	mux.Handle("/", websocket.Handler(mod.Handle))
	//一个测试html
	mux.HandleFunc("/web", WebSocketHTMLHandlego)
	//你也可以在外面继续扩展

	mod.httpServer.Handler = mux
}

//Start IModule   接口实现
func (mod *WebSocketModule) Start() {

	//启动的协程
	mod.threadgo.Go(func(ctx context.Context) {
		mod.wg.Add(1)
		defer mod.wg.Done()
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

//Stop IModule 接口实现
func (mod *WebSocketModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		Logger.PError(err, "Close websocket Module:")
	}
	mod.threadgo.CloseWait()
	// mod.wsmaplock.Lock()
	// m := mod.wsmap
	// for k := range m {
	// 	k.Close()
	// }
	// mod.wsmaplock.Unlock()
	// mod.wg.Wait()
	Logger.PStatus("websocket Module Stop.")
}

//PrintStatus IModule 接口实现，打印状态
func (mod *WebSocketModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n		websocket Module    :%d/%d/%d	(connum/get/sendnum)",
		len(mod.wsmap),
		atomic.AddInt64(&mod.getnum, 0),
		atomic.AddInt64(&mod.sendnum, 0))
}

//Handle http发来的所有请求都会到这个方法来
func (mod *WebSocketModule) Handle(conn *websocket.Conn) {
	//连接加入管理器，好用来关闭
	mod.mapadd(conn)
	defer mod.mapdel(conn)

	//标注子连接是不是都停下来
	mod.wg.Add(1)
	defer mod.wg.Done()
	defer conn.Close()

	//用来管理连接下开的子协程
	runobj := threads.NewThreadGoByGo(mod.threadgo)
	request := make([]byte, 10240)
	// msgbuff := make([]byte, 0, 20480)
	//发给下面的连接对象，可以自定义一些信息和回调
	wsconn := new(event.WebSocketModel)
	wsconn.Conn = conn
	wsconn.KeyID = -1
	wsname := conn.Request().RemoteAddr
	if mod.WebSocketOnlineFun != nil {
		wsname = mod.WebSocketOnlineFun(conn)
	}
	//发消息来说明这个用户掉线了
	defer func() {
		Logger.PInfoKey("%s websocket client closeing.", wsconn.KeyID, wsname)
		runobj.CloseWait() //要等下面的逻辑都处理完了，才可以运行下面的代码，保证保存的逻辑
		//用来处理发生连接关闭的时候，要处理的事
		if wsconn.CloseFun != nil {
			wsconn.CloseFun(wsconn)
		}
		Logger.PInfoKey("%s websocket client close.", wsconn.KeyID, wsname)
	}()
	Logger.PInfoKey("%s websocket client open!", wsconn.KeyID, wsname)
	runchan := make(chan bool, 8) //用来处理超时
	mod.threadgo.Go(
		func(ctx context.Context) {
			timeout := time.NewTimer(time.Duration(mod.cg.Timeout) * time.Second)
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
						timeout.Reset(time.Duration(mod.cg.Timeout) * time.Second)
					} else {
						return
					}
				}
			}

			//超时关连接
		})
	runobj.Try(
		func(ctx context.Context) {
		listen:
			for {
				readLen, err := conn.Read(request)
				if err != nil {
					if err == io.EOF {
						runchan <- false
						//fmt.Println("客户端断开链接，")
						break listen
					} else {
						//Logger.PErrorKey(err, "websocket error.", wsconn.KeyID)
						// fmt.Println(err)
					}
				}
				Logger.PInfoKey(string(request[:readLen]), wsconn.KeyID)
				etjs := make(event.JsonMap)
				err = json.Unmarshal(request[:readLen], &etjs)
				if err != nil {
					break listen
				}
				atomic.AddInt64(&mod.getnum, 1)
				code := etjs.GetAction() //["ACTION"]) //可能会出错，不知道有没有捕获
				call := mod.RouteFun(code)
				if call == nil {
					Logger.PInfoKey("nothing action:%d!", wsconn.KeyID, code)
				} else {
					runchan <- true
					call(etjs, wsconn, runobj) //调用委托的消息处理方法
				}
			}
		},
		nil,
		nil,
	)

	// msgbuff = append(msgbuff, request[:readLen]...)
	// for {
	// 	//也有可能一次收到好几条消息
	// 	msglen := int(binary.BigEndian.Uint16(msgbuff))
	// 	if msglen > mod.cg.MsgMaxLen {
	// 		//消息太长
	// 		break listen
	// 	}
	// 	//消息完整了，就可以开始处理
	// 	if len(msgbuff) < msglen {
	// 		continue listen
	// 	}
	// 	etjs := make(map[string]interface{})
	// 	err = json.Unmarshal(msgbuff[2:msglen], &etjs)
	// 	if err != nil {
	// 		break listen
	// 	}
	// 	//已解出来消息要从缓存中去掉
	// 	msgbuff = append(msgbuff[msglen:])

	// 	atomic.AddInt64(&mod.getnum, 1)
	// 	code := util.Convert.ToInt32(etjs["ACTION"]) //可能会出错，不知道有没有捕获
	// 	call := mod.RouteFun(code)
	// 	if call == nil {
	// 		Logger.PInfo("nothing action:%d!", code)
	// 	} else {
	// 		call(etjs, conn, runobj) //调用委托好的消息处理方法
	// 	}
	// }

}

//WebSocketHTMLHandlego 默认的所有没定义的处理请求
func WebSocketHTMLHandlego(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "web.html")
}

func (mod *WebSocketModule) mapadd(conn *websocket.Conn) {
	mod.wsmaplock.Lock()
	defer mod.wsmaplock.Unlock()
	if mod.wsmap == nil {
		conn.Close()
	} else {
		mod.wsmap[conn] = true
	}
}
func (mod *WebSocketModule) mapdel(conn *websocket.Conn) {
	mod.wsmaplock.Lock()
	defer mod.wsmaplock.Unlock()
	if mod.wsmap != nil {
		delete(mod.wsmap, conn)
	}
}

//GetPlayerNum用户连接数量
func (mod *WebSocketModule) GetPlayerNum() int {
	return len(mod.wsmap)
}
