package module_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/event"
	"github.com/buguang01/bige/module"
)

func TestModule(t *testing.T) {
	Logger.Init(0, "logs", Logger.LogModeFmt)
	m := module.NewHTTPModule(&module.HTTPConfig{
		HTTPAddr: ":8080",
		Timeout:  10,
	})

	m.RouteFun = GetEventByData //设置路由
	// m.TimeoutFun = TimeoutCallback //设置超时

	m.Init()

	m.Start()

	time.Sleep(600 * time.Second)
	m.Stop()

	Logger.LogClose()

}

//GetEventByData 路由器的例子
func GetEventByData(code int) event.HTTPcall {
	switch code {
	case 1001:
		return HTTPGetMsgHandle
	case 1002:
		return HTTPGetMsgHandle2
	}
	return nil
}

func HTTPGetMsgHandle(et event.JsonMap, w http.ResponseWriter) {

	//这是一段例子，也是我自己定义的标准回复信息
	jsuser := make(event.JsonMap)
	jsuser["Member"] = MemberMD{100001, "xiacs11111"}
	// resultjs["ACTIONKEY"] = et.ActionKey

	event.HTTPReplyMsg(w, et, 0, jsuser)
	fmt.Println("timeout run")

	//用来测试消息处理超时了会怎么样
}

func HTTPGetMsgHandle2(et event.JsonMap, w http.ResponseWriter) {

	//这是一段例子，也是我自己定义的标准回复信息
	jsuser := make(event.JsonMap)
	jsuser["Member"] = MemberMD{100002, "xiacs222222"}
	// resultjs["ACTIONKEY"] = et.ActionKey

	time.Sleep(40 * time.Second)
	event.HTTPReplyMsg(w, et, 0, jsuser)
	fmt.Println("timeout run")

	//用来测试消息处理超时了会怎么样
}

type Aclass struct {
	atext string
}
type Bclass struct {
	Aclass
	btext string
}

func TestTmp(t *testing.T) {
	var b interface{} = new(Bclass)
	a, ok := b.(*Aclass)
	if ok {

	}
	fmt.Print(a)
	fmt.Print(b)
}

func TestTime(t *testing.T) {
	// g := threads.NewThreadGo()
	m := make(map[int]*event.SqlDataModel, 10240000)

	for i := 0; i < 10240000; i++ {
		m[i] = new(event.SqlDataModel)
		k := rand.Intn(100)
		if k < 10 {
			delete(m, rand.Intn(i))
		}
		// g.Go(func(ctx context.Context) {
		// 	<-ctx.Done()
		// })
	}
	fmt.Println(len(m))
	// fmt.Println("end1")
	time.Sleep(10 * time.Second)
	for i := 0; i < 10240000; i++ {
		m[i] = new(event.SqlDataModel)
		k := rand.Intn(100)
		if k < 10 {
			delete(m, rand.Intn(i))
		}
		// g.Go(func(ctx context.Context) {
		// 	<-ctx.Done()
		// })
	}
	// fmt.Println(len(m))
	fmt.Println("end2")

}
