package module

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/event"
	"github.com/buguang01/bige/json"
	"github.com/buguang01/bige/threads"
)

//HTTPConfig httpmodule的配置
type HTTPConfig struct {
	HTTPAddr string //监听地址
	Timeout  int
}

//HTTPModule ...
//http连接模块
type HTTPModule struct {
	HTTPAddr   string         //HTTP监听的地址
	httpServer *http.Server   //HTTP请求的对象
	cg         *HTTPConfig    //从配置表中读进来的数据
	wg         sync.WaitGroup //用来确定是不是关闭了
	getnum     int64          //收到的总消息数
	runing     int64          //当前在处理的消息数
	// failnum    int64          //发生问题的消息数

	RouteFun   func(code int) event.HTTPcall                         //用来生成事件处理器的工厂
	TimeoutFun event.HTTPcall                                        //超时时的回调方法
	GetIPFun   func(w http.ResponseWriter, req *http.Request) string //拿IP的方法
}

//NewHTTPModule 生成一个新的HTTP的对象
func NewHTTPModule(configmd *HTTPConfig) *HTTPModule {
	result := &HTTPModule{
		cg:       configmd,
		HTTPAddr: configmd.HTTPAddr,
	}

	return result
}

//Init IModule接口的实现
func (mod *HTTPModule) Init() {

	mod.httpServer = &http.Server{
		Addr:         mod.cg.HTTPAddr,
		WriteTimeout: time.Duration(mod.cg.Timeout) * time.Second,
	}
	//还可以加别的参数，已后再加，有需要再加
	mux := http.NewServeMux()
	//这个是主要的逻辑
	mux.HandleFunc("/request", mod.Handle)
	//这个只是防止404用的
	mux.HandleFunc("/", NullHandle)
	mux.HandleFunc("/web", HTMLHandlego)
	//你也可以在外面继续扩展

	mod.httpServer.Handler = mux

	mod.TimeoutFun = TimeoutRun
}

//Start IModule   接口实现
func (mod *HTTPModule) Start() {

	//启动的协程
	go func() {
		mod.wg.Add(1)
		defer mod.wg.Done()
		Logger.PStatus("HTTP Module Start!")
		err := mod.httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				Logger.PStatus("Http run Server closed under requeset!!")
				// log.Print("Server closed under requeset!!")
			} else {
				Logger.PFatal("Server closed unexpecteed:" + err.Error())
				// log.Fatal("Server closed unexpecteed!!")
			}
		}
	}()
}

//Stop IModule 接口实现
func (mod *HTTPModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		Logger.PError(err, "Close HttpModule:")
	}
	mod.wg.Wait()
	Logger.PStatus("HTTP Module Stop")
}

//PrintStatus IModule 接口实现，打印状态
func (mod *HTTPModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n		HTTP Module         :%d/%d		(get/runing)",
		atomic.AddInt64(&mod.getnum, 0),
		atomic.AddInt64(&mod.runing, 0))
}

//Handle http发来的所有请求都会到这个方法来
func (mod *HTTPModule) Handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mod.wg.Add(1)
	defer mod.wg.Done()
	atomic.AddInt64(&mod.getnum, 1)
	atomic.AddInt64(&mod.runing, 1)
	defer atomic.AddInt64(&mod.runing, -1)
	timeout := time.NewTimer(mod.httpServer.WriteTimeout - 2*time.Second)
	request := req.FormValue("json")
	etjs := make(event.JsonMap)
	err := json.Unmarshal([]byte(request), &etjs)
	if err != nil {
		w.Write([]byte("json error."))
		return
	}
	Logger.PInfo(request)
	threads.Try(
		func() {
			ip := req.RemoteAddr
			if mod.GetIPFun != nil {
				ip = mod.GetIPFun(w, req)
			}
			action := etjs.GetAction()
			call := mod.RouteFun(action)
			etjs["IP"] = ip
			if call == nil {
				Logger.PInfo("nothing action:%d!", action)
				w.Write([]byte("nothing action"))
			} else {
				g := threads.NewGoRun(
					func() {
						call(etjs, w)
					},
					nil)
				select {
				case <-g.Chanresult:
					timeout.Stop()
					//上面那个运行完了
					break
				case <-timeout.C:
					//上面那个可能还没有运行完，但是超时了要返回了
					Logger.PDebug("http time msg:%s", request)
					if mod.TimeoutFun != nil {
						mod.TimeoutFun(etjs, w)
					}
					break
				}
				//调用委托好的消息处理方法
			}
		},
		func(err interface{}) {
			Logger.PFatal(err)
			//如果出异常了，跑这里
			w.Write([]byte("catch!"))
		},
		nil)

	// w.Write(mod.TimeoutFun(etdata))

}

//NullHandle 默认的所有没定义的处理请求
func NullHandle(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world!"))
}

//HTMLHandlego 默认的所有没定义的处理请求
func HTMLHandlego(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "web.html")
}

//TimeoutRun 默认的超时调用
func TimeoutRun(et event.JsonMap, w http.ResponseWriter) {
	event.HTTPReplyMsg(w, et, -1, nil)
}
