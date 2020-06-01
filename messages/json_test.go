package messages_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/buguang01/bige/event"
	"github.com/buguang01/bige/messages"
)

func TestJson(t *testing.T) {
	handle := messages.JsonMessageHandleNew()
	handle.SetRoute(1, make(event.JsonMap))
	et := make(event.JsonMap)
	et["abc"] = 1
	et["qqqq"] = "abcdadfa"
	databuff, err := handle.Marshal(1, et)
	if err == nil {
		fmt.Printf("databuff:%v", databuff)
	} else {
		fmt.Println(err)
		return
	}
	result, err := handle.Unmarshal(databuff)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("result:%v", result)
	}

}

type TestMsg struct {
	messages.Message
	// ActionID uint32 `json:"ACTIONID"`
	MemberID int `json:"MEMBERID"`
}

func TestMsgfmt(t *testing.T) {
	msg := new(TestMsg)
	msg.ActionID = 1001
	msg.MemberID = 10020202
	b, _ := json.Marshal(msg)
	fmt.Println(string(b))
	msg2 := new(TestMsg)
	json.Unmarshal(b, msg2)
	fmt.Printf("%+v", msg2)

}
