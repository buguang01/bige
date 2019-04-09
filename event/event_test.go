package event_test

import (
	"buguang01/gsframe/event"
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonEvent(t *testing.T) {

	et := &EventUserLogin{}

	p := func(es event.IMsgEvent) {
		json.Unmarshal([]byte(`
		{"action":1000,"appt":666}
		`), es)
		fmt.Println(es)
		es.WorkHandle()
	}
	p(et)
	fmt.Println(et)
	et.WorkHandle()

}

//子类也可以正常使用
type EventUserLogin struct {
	event.MsgEventBase
	Appt int
}

// func (et *EventUserLogin) WorkHandle() {
// 	fmt.Printf("%T 实现WorkHandle", et)

// }
