package modules_test

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/buguang01/bige/modules"
	"github.com/buguang01/util/threads"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/event"
	"github.com/buguang01/bige/messages"
)

var (
	WebmodulesEx *modules.WebModule
)

func TestWeb(t *testing.T) {

	Logger.Init(0, "logs", Logger.LogModeFmt)
	WebmodulesEx := modules.NewWebModule()
	WebmodulesEx.RouteHandle = messages.HttpJsonMessageHandleNew()
	action := Msgone{}.GetAction()
	WebmodulesEx.RouteHandle.SetRoute(action, &Msgone{})

	// m.TimeoutFun = TimeoutCallback //设置超时

	WebmodulesEx.Init()

	WebmodulesEx.Start()

	time.Sleep(600 * time.Second)
	WebmodulesEx.Stop()

	Logger.LogClose()
}

type Msgone struct {
	UserName string
	PassWord string
	MemberID int
}

func (msg Msgone) GetAction() uint32 {
	return 1001
}
func (msg *Msgone) HttpDirectCall(w http.ResponseWriter, req *http.Request) {
	jsuser := make(event.JsonMap)
	jsuser["Name"] = msg.UserName
	jsuser["ACTION"] = int(msg.GetAction())
	event.HTTPReplyMsg(w, jsuser, 0, jsuser)
}

func TestCtx(t *testing.T) {
	thgo := threads.NewThreadGo()
	thgo.CloseWait()
	for {
		select {
		case <-thgo.Ctx.Done():
			fmt.Println("ctx")
		default:
			fmt.Println("def")
		}
		time.Sleep(time.Second)
	}

}

func TestAddInt(t *testing.T) {
	// ch := make(chan int, 8)
	// ch <- 1
	// ch <- 2
	// ch <- 3
	// close(ch)
	// tk := time.NewTimer(time.Second * 10)
	// for {
	// 	select {
	// 	case <-ch:
	// 		fmt.Println("tk.c then")
	// 		time.Sleep(time.Second)
	// 	case <-tk.C:
	// 		fmt.Println("tk.c ")
	// 		time.Sleep(time.Second)

	// 	}
	// }

	var i int64 = 0
	addint(&i)
	fmt.Println(i)
}

func addint(i *int64) {
	atomic.AddInt64(i, 1)
	defer atomic.AddInt64(i, 10)
	atomic.AddInt64(i, 100)
}
