package module

import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//HTTPConfig httpmodule的配置
type HTTPConfig struct {
	HTTPAddr string //监听地址
	Timeout  int32
}

//HTTPModule ...
//http连接模块
type HTTPModule struct {
	HTTPAddr   string         //HTTP监听的地址
	httpServer *http.Server   //HTTP请求的对象
	cg         *HTTPConfig    //从配置表中读进来的数据
	wg         sync.WaitGroup //用来确定是不是关闭了

	RouteFun   func(code int32) event.IHTTPMsgEVent //用来生成事件处理器的工厂
	TimeoutFun func(et event.IHTTPMsgEVent) []byte  //超时时的回调方法
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
func (mod *HTTPModule) Init() error {
	// if configpath == "" {
	// 	configpath = "config/http.json"
	// }
	// filedb, err := ioutil.ReadFile(configpath)
	// if err != nil {
	// 	//有问题就要退出的
	// 	loglogic.PFatal(err)
	// 	return err
	// }
	// if err = json.Unmarshal([]byte(filedb), &mod.cg); err != nil {
	// 	//有问题就要退出
	// 	loglogic.PFatal(err)
	// 	return err
	// }

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
	//你也可以在外面继续扩展

	mod.httpServer.Handler = mux
	return nil
}

//Start IModule   接口实现
func (mod *HTTPModule) Start() {

	//启动的协程
	go func() {
		mod.wg.Add(1)
		defer mod.wg.Done()
		err := mod.httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				loglogic.PStatus("Http run Server closed under requeset!!")
				// log.Print("Server closed under requeset!!")
			} else {
				loglogic.PFatal("Server closed unexpecteed:" + err.Error())
				// log.Fatal("Server closed unexpecteed!!")
			}
		}
	}()
}

//Stop IModule 接口实现
func (mod *HTTPModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		loglogic.PError("Close HttpModule:" + err.Error())
	}
	mod.wg.Wait()
	loglogic.PStatus("Http close")
}

//Handle http发来的所有请求都会到这个方法来
func (mod *HTTPModule) Handle(w http.ResponseWriter, req *http.Request) {
	timeout := time.NewTicker(mod.httpServer.WriteTimeout - 2*time.Second)
	v := req.FormValue("action")
	if v == "" {
		w.Write([]byte("nothing action!!!"))
		return
	}
	action, _ := strconv.Atoi(v)
	loglogic.PDebug(req.URL.String())

	etdata := mod.RouteFun(int32(action))
	if etdata == nil {
		//没拿到就是没有定义过这个消息
		loglogic.PStatus("Undefined action:%d", action)
		w.Write([]byte(fmt.Sprintf("Undefined action:%d", action)))
		return
	}
	dval := req.FormValue("jsdata")
	if dval == "" {
		w.Write([]byte("nothing jsdata!!!"))
		return
	}
	if err := json.Unmarshal([]byte(dval), etdata); err != nil {
		loglogic.PStatus("json Unmarshal error:%v", err)
		w.Write([]byte("json Unmarshal fail!"))
		return
	}
	if int32(action) != etdata.(event.IMsgEvent).GetAction() {
		loglogic.PStatus("action error Url:%s", req.URL)
		w.Write([]byte("action error!"))
		return

	}
	resultchan := etdata.HTTPGetMsgHandle()
	select {
	case d := <-resultchan:
		w.Write(d)

	case <-timeout.C:
		w.Write(mod.TimeoutFun(etdata))

	}
}

//NullHandle 默认的所有没定义的处理请求
func NullHandle(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world!"))
}
