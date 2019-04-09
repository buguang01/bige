package module_test

import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/module"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestModule(t *testing.T) {
	loglogic.Init(0, "logs")
	m := module.NewHTTPModule(&module.HTTPConfig{
		HTTPAddr: ":8080",
		Timeout:  10,
	})

	m.RouteFun = GetEventByData    //设置路由
	m.TimeoutFun = TimeoutCallback //设置超时

	m.Init()

	m.Start()

	time.Sleep(600 * time.Second)
	m.Stop()

	loglogic.LogClose()

}

//GetEventByData 路由器的例子
func GetEventByData(code int32) event.IHTTPMsgEVent {
	switch code {
	case 1001:
		return &event.MsgEventBase{}
	case 1002:
		return &MsgEventTimeoutExample{}
	}
	return nil
}

//TimeoutCallback 超时例子
func TimeoutCallback(etdata event.IHTTPMsgEVent) []byte {
	et := etdata.(event.IMsgEvent)
	resultjs := make(map[string]interface{})
	resultjs["ACTION"] = et.GetAction()
	resultjs["ACTIONCOM"] = -1
	resultjs["ACTIONKEY"] = et.GetActionKey()
	resultjs["GETMSG"] = et
	resultb, _ := json.Marshal(resultjs)
	return resultb
}

type MsgEventTimeoutExample struct {
	event.MsgEventBase
}

func (et *MsgEventTimeoutExample) HTTPGetMsgHandle() <-chan []byte {
	result := make(chan []byte, 1)

	//这是一段例子，也是我自己定义的标准回复信息
	resultjs := make(map[string]interface{})
	resultjs["ACTION"] = et.Action
	resultjs["ACTIONCOM"] = 0
	resultjs["ACTIONKEY"] = et.ActionKey
	resultb, _ := json.Marshal(resultjs)
	go func() {
		time.Sleep(40 * time.Second)
		result <- resultb
		close(result)
		fmt.Println("timeout run")
	}()
	return result
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
