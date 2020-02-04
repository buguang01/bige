package modules_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/buguang01/bige/modules"

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
