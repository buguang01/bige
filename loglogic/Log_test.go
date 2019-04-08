package loglogic

import (
	"testing"
)

func TestLog(t *testing.T) {
	defer LogClose()

	defer func() {
		r := recover()
		if r != nil {
			PError(r)

		}
	}()
	Init(20, "logs")
	SetListenKeyID(1001)
	PDebug("test1")
	PInfo("test2")
	PInfoKey("test3", 1001)
	PDebugKey("test4", 1002)
	PDebugKey("test5", 1001)
	//errs := fmt.Errorf(string(debug.Stack()))
	panic("panicinfo")

}
