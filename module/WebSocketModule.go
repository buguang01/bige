package module

import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/threads"
	"buguang01/gsframe/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

//如果一个服务器，又要接受客户端的，又要接受其他服务器的消息转发
//建议开二个这样的服务器，然后用一个BridgeModule,做消息的转发
//

//WebSocketConfig websocket的配置
type WebSocketConfig struct {
	Addr      string //websocket 监听地址
	Timeout   int32  //超时时间 （秒）
	MsgMaxLen int    //一个消息最大长度B

}

//WebSocketModule socket监听模块
type WebSocketModule struct {
	Addr       string           //HTTP监听的地址
	httpServer *http.Server     //HTTP请求的对象
	cg         *WebSocketConfig //从配置表中读进来的数据
	wg         sync.WaitGroup   //用来确定是不是关闭了
	getnum     int64            //收到的总消息数
	sendnum    int64            //当前在处理的消息数

	wsmap     map[*websocket.Conn]bool //所有的连接
	wsmaplock sync.Mutex               //上面那个对象的锁

	RouteFun func(code int32) event.WebSocketCall //用来生成事件处理器的工厂

	// wslist         map[uint32]*WsConnModel    //websocket 列表
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
	// result.wslist = make(map[uint32]*WsConnModel)
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
	go func() {
		mod.wg.Add(1)
		defer mod.wg.Done()
		loglogic.PStatus("websocket Module Start!")
		err := mod.httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				loglogic.PStatus("websocket run Server closed under requeset!!")
				// log.Print("Server closed under requeset!!")
			} else {
				loglogic.PFatal("Server closed unexpecteed:" + err.Error())
				// log.Fatal("Server closed unexpecteed!!")
			}
		}
	}()
}

//Stop IModule 接口实现
func (mod *WebSocketModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		loglogic.PError("Close websocket Module:" + err.Error())
	}
	mod.wsmaplock.Lock()
	m := mod.wsmap
	for k := range m {
		k.Close()
	}
	mod.wsmaplock.Unlock()
	mod.wg.Wait()
	loglogic.PStatus("websocket Module Stop.")
}

//PrintStatus IModule 接口实现，打印状态
func (mod *WebSocketModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\twebsocket Module:	%d/%d	(get/sendnum)",
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
	runobj := new(threads.ThreadGo)
	request := make([]byte, 10240)
	// msgbuff := make([]byte, 0, 20480)
	//发给下面的连接对象，可以自定义一些信息和回调
	wsconn := new(event.WebSocketModel)
	wsconn.Conn = conn

	//发消息来说明这个用户掉线了
	defer func() {
		loglogic.PInfo("websocket client closeing.")
		runobj.Wg.Wait() //要等下面的逻辑都处理完了，才可以运行下面的代码，保证保存的逻辑
		//用来处理发生连接关闭的时候，要处理的事
		if wsconn.CloseFun != nil {
			wsconn.CloseFun(wsconn)
		}
		loglogic.PInfo("websocket client close.")
	}()
	loglogic.PInfo("websocket client open!")
	runchan := make(chan bool, 8) //用来处理超时
	threads.GoTry(
		func() {
			timeout := time.NewTicker(time.Duration(mod.cg.Timeout) * time.Second)
			for {
				select {
				case <-timeout.C:
					conn.Close()
				case ok := <-runchan:
					if ok {
						timeout = time.NewTicker(time.Duration(mod.cg.Timeout) * time.Second)
					} else {
						return
					}
				}
			}

			//超时关连接
		}, nil, nil)
listen:
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			if err == io.EOF {
				runchan <- false
				//fmt.Println("客户端断开链接，")
				break listen
			} else {
				fmt.Println(err)
			}
		}
		loglogic.PInfo(string(request[:readLen]))
		etjs := make(event.JsonMap)
		err = json.Unmarshal(request[:readLen], &etjs)
		if err != nil {
			break listen
		}
		atomic.AddInt64(&mod.getnum, 1)
		code := util.Convert.ToInt32(etjs["ACTION"]) //可能会出错，不知道有没有捕获
		call := mod.RouteFun(code)
		if call == nil {
			loglogic.PInfo("nothing action:%d!", code)
		} else {
			runchan <- true
			call(etjs, wsconn, runobj) //调用委托的消息处理方法
		}

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
		// 		loglogic.PInfo("nothing action:%d!", code)
		// 	} else {
		// 		call(etjs, conn, runobj) //调用委托好的消息处理方法
		// 	}
		// }
	}
}

//WebSocketHTMLHandlego 默认的所有没定义的处理请求
func WebSocketHTMLHandlego(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "web.html")
}

// //RegisterSocket 注册SOCKRT
// func (mod *WebSocketModule) RegisterSocket(conn *websocket.Conn, hash uint32, f func()) {
// 	mod.wslock.Lock()
// 	defer mod.wslock.Unlock()
// 	mod.wslist[hash] = conn
// 	mod.wsclosefunlist[conn] = f
// }

// //RemoveSocket 移除连接
// func (mod *WebSocketModule) RemoveSocket(conn *websocket.Conn, hash uint32) {
// 	mod.wslock.Lock()
// 	defer mod.wslock.Unlock()
// 	delete(mod.wslist, hash)
// 	delete(mod.wsclosefunlist, conn)
// }

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
