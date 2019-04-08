package module

import (
	"buguang01/gsframe/config"
	"buguang01/gsframe/loglogic"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//HTTPModule ...
//http连接模块
type HTTPModule struct {
	HTTPAddr   string
	httpServer *http.Server
	wgmd       sync.WaitGroup
	quitch     chan int
	cg         *config.HTTPConfig
}

//Init IModule接口的实现
func (mod *HTTPModule) Init(configpath string) error {
	if configpath == "" {
		configpath = "config/http.json"
	}
	filedb, err := ioutil.ReadFile(configpath)
	if err != nil {
		//有问题就要退出的
		loglogic.PFatal(err)
		return err
	}
	if err = json.Unmarshal([]byte(filedb), &mod.cg); err != nil {
		//有问题就要退出
		loglogic.PFatal(err)
		return err
	}

	mod.httpServer = &http.Server{
		Addr:         mod.cg.HTTPAddr,
		WriteTimeout: mod.cg.Timeout * time.Second,
	}
	//还可以加别的参数，已后再加，有需要再加

	mod.quitch = make(chan int)
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handle)
	mod.httpServer.Handler = mux
	return nil
}

//Start IModule   接口实现
func (mod *HTTPModule) Start() {

	//启动的协程
	go func() {
		mod.wgmd.Add(1)
		defer mod.wgmd.Done()
		err := mod.httpServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				loglogic.PStatus("Server closed under requeset!!")
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
	mod.wgmd.Wait()

}

//Handle http发来的所有请求都会到这个方法来
func Handle(w http.ResponseWriter, req *http.Request) {
	// w.Write([]byte("Hello"))

}
