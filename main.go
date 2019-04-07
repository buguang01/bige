package main

import (
	"buguang01/gsframe/loglogic"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("test")
	loglogic.Init(0, "logs")
	loglogic.PDebug("test1")
	loglogic.PDebug("test2")
	time.Sleep(60 * time.Second)
}

func main1() {

	server := &http.Server{
		Addr:         ":8001",
		WriteTimeout: 4 * time.Second,
	}
	quit := make(chan os.Signal)
	//创建chan，用来指示我要退出这个服务器了，麻烦帮忙关闭一下
	signal.Notify(quit, os.Interrupt)
	//注册这个通知事件，一旦受到这个singal，发送一个对象到这个chan当中，当我接收到任意对象之后，我就知道服务器该退出了
	mux := http.NewServeMux()
	//自己创建servemux，然后使用自己的handle方法,mux就是实现了handler接口的一个变量
	// mux.Handle("/", &myHandler{})
	mux.HandleFunc("/", SayHolle)
	//默认的mux中根路由包含了所有的未匹配的路由
	mux.HandleFunc("/bye", SayBye)
	//将mux集成到server当中，server.Handle也是handle类型的接口，所以可以直接赋值
	server.Handler = mux
	//创建一个gorouting 专门接收这个chan
	go func() {
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal("Close server:", err)
		}
	}()

	log.Println("Start version v3")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Print("Server closed under requeset!!")
		} else {
			log.Fatal("Server closed unexpecteed!!")
		}

	}
	log.Println("Server exit!!")
	time.Sleep(10 * time.Second)
	log.Println("Server exit!!")

}

type myHandler struct{} //自己定义handler结构
//实现myHandler的ServeHTPP方法
func (my *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello this is version 3!,the requeset URL is:" + r.URL.String()))
	//这里可以打印出完整的URL，响应的都是根路由
}
func SayHolle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bye bye, this is version v11111111"))

}

func SayBye(w http.ResponseWriter, r *http.Request) {
	log.Println("Saybye start")

	time.Sleep(3 * time.Second)

	w.Write([]byte("Bye bye, this is version v3"))
	log.Println("Saybye end")

	//进行一个流式传递，将字符串转换为byte类型
}
