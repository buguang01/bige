package module

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"bug.guang/config"
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
func (mod *HTTPModule) Init(configpath string) {
	if configpath == "" {
		configpath = "config/http.json"
	}
	filedb, err := ioutil.ReadFile(configpath)
	if err != nil {
		//有问题就要退出的

	}
	if err = json.Unmarshal([]byte(filedb), &mod.cg); err != nil {
		//有问题就要退出
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
				// log.Print("Server closed under requeset!!")
			} else {
				// log.Fatal("Server closed unexpecteed!!")
			}

		}
	}()
}

//Stop IModule 接口实现
func (mod *HTTPModule) Stop() {
	if err := mod.httpServer.Close(); err != nil {
		log.Fatal("Close HttpModule:", err)
	}
	mod.wgmd.Wait()

}

//Handle http发来的所有请求都会到这个方法来
func Handle(w http.ResponseWriter, req *http.Request) {
	// w.Write([]byte("Hello"))

}
