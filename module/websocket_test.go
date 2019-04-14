package module_test

import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/module"
	"buguang01/gsframe/threads"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	loglogic.Init(0, "logs")
	m := module.NewWSModule(&module.WebSocketConfig{
		Addr:      ":8080",
		Timeout:   10,
		MsgMaxLen: 10240,
	})

	m.RouteFun = RouteFun

	m.Init()

	m.Start()

	time.Sleep(30 * time.Second)
	m.Stop()

	loglogic.LogClose()
}

func RouteFun(code int32) event.WebSocketCall {
	if code == 1001 {
		return func(et event.JsonMap, wsmd *event.WebSocketModel, runobj *threads.ThreadGo) {
			//在新线程上跑
			runobj.Go(func() {
				time.Sleep(30 * time.Second)
				jsuser := make(event.JsonMap)
				jsuser["Member"] = MemberMD{100001, "xiacs"}
				event.WebSocketReplyMsg(wsmd, et, 0, jsuser)
			})
		}
	} else if code == 1002 {
		return func(et event.JsonMap, wsmd *event.WebSocketModel, runobj *threads.ThreadGo) {
			//在当前线程上跑
			runobj.Try(func() {
				// time.Sleep(10 * time.Second)
				event.WebSocketReplyMsg(wsmd, et, 0, nil)
			}, nil, nil)
		}
	}
	return nil
}

type MemberMD struct {
	MemberID int32
	UserName string
}
