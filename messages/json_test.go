package messages_test

import (
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
