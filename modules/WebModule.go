package modules

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util/threads"
)

//设置Web地址
func WebModuleSetIpPort(ipPort string) options {
	return func(mod IModule) {
		mod.(*WebModule).ipPort = ipPort
	}
}

//设置超时时间（秒）
func WebModuleSetTimeout(timeout time.Duration) options {
	return func(mod IModule) {
		mod.(*WebModule).timeout = timeout * time.Second
	}
}

//设置超时回调方法
func WebModuleSetTimeoutFunc(timeoutfunc func(webmsg messages.IHttpMessageHandle,
	w http.ResponseWriter, req *http.Request)) options {
	return func(mod IModule) {
		mod.(*WebModule).timeoutFun = timeoutfunc
	}
}

type WebModule struct {
	ipPort      string                  //监听地址
	timeout     time.Duration           //超时时时
	RouteHandle messages.IMessageHandle //消息路由
	timeoutFun  func(webmsg messages.IHttpMessageHandle,
		w http.ResponseWriter, req *http.Request) //超时回调，把超时的消息传入
	getnum     int64             //收到的总消息数
	runing     int64             //当前在处理的消息数
	httpServer *http.Server      //HTTP请求的对象
	thgo       *threads.ThreadGo //协程管理器
}

func NewWebModule(opts ...options) *WebModule {
	result := &WebModule{
		ipPort:      ":8080",
		timeout:     30 * time.Second,
		timeoutFun:  webTimeoutRun,
		getnum:      0,
		runing:      0,
		thgo:        threads.NewThreadGo(),
		RouteHandle: messages.JsonMessageHandleNew(),
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

//Init 初始化
func (mod *WebModule) Init() {
	mod.httpServer = &http.Server{
		Addr:         mod.ipPort,
		WriteTimeout: mod.timeout,
	}
	//还可以加别的参数，已后再加，有需要再加
	mux := http.NewServeMux()
	//这个是主要的逻辑
	mux.HandleFunc("/", mod.Handle)
	mod.httpServer.Handler = mux
}

//Start 启动
func (mod *WebModule) Start() {
	mod.thgo.Go(func(ctx context.Context) {
		Logger.PStatus("Web Module Start!")
		err := mod.httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				Logger.PStatus("Web run Server closed under requeset!!")
			} else {
				Logger.PError(err, "Server closed unexpecteed.")
			}
		}
	})
}

//Stop 停止
func (mod *WebModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		Logger.PError(err, "Close Web Module.")
	}
	mod.thgo.CloseWait()
	Logger.PStatus("Web Module Stop.")
}

//PrintStatus 打印状态
func (mod *WebModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tWeb Module\t:%d/%d\t(get/runing)",
		atomic.LoadInt64(&mod.getnum),
		atomic.LoadInt64(&mod.runing))
}

//Handle http发来的所有请求都会到这个方法来
func (mod *WebModule) Handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mod.thgo.Wg.Add(1)
	defer mod.thgo.Wg.Done()
	atomic.AddInt64(&mod.getnum, 1)
	atomic.AddInt64(&mod.runing, 1)
	defer atomic.AddInt64(&mod.runing, -1)
	timeout := time.NewTimer(mod.httpServer.WriteTimeout - 2*time.Second)
	buff, _ := ioutil.ReadAll(req.Body)
	// fmt.Println(string(buff))
	msg, err := mod.RouteHandle.Unmarshal(buff)
	if err != nil {
		Logger.PInfo("RouteHandle Unmarshal Error:%s", err.Error())
		return
	}
	modmsg, ok := msg.(messages.IHttpMessageHandle)
	if !ok {
		Logger.PInfo("Not is Web Msg:%+v", msg)
		return
	} else {
		Logger.PInfo("Web Get Msg:%+v", msg)
	}

	threads.Try(
		func() {
			g := threads.NewGoRun(
				func() {
					modmsg.HttpDirectCall(w, req)
				},
				nil)
			select {
			case <-g.Chanresult:
				timeout.Stop()
				//上面那个运行完了
				break
			case <-timeout.C:
				//上面那个可能还没有运行完，但是超时了要返回了
				Logger.PDebug("web timeout msg:%+v", modmsg)
				if mod.timeoutFun != nil {
					mod.timeoutFun(modmsg, w, req)
				}
				break
			}
			//调用委托好的消息处理方法
		},
		func(err interface{}) {
			Logger.PFatal(err)
			//如果出异常了，跑这里
			w.Write([]byte("catch!"))
		},
		nil)
}

func webTimeoutRun(webmsg messages.IHttpMessageHandle, w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("timeout Run!"))
}
