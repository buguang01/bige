package event_test

import (
	"github.com/buguang01/gsframe/event"
	"fmt"
	"testing"
)

func TestMemory(t *testing.T) {
	var tmp event.IMemoryModel = new(TmpMemory)

	tmp.Handle(nil)
}

type TmpMemory struct {
	event.MemoryModel
	Num int
}

func (this *TmpMemory) Run() {
	fmt.Println("tmpmemory")
}
