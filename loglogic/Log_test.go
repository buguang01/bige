package loglogic_test

import (
	"buguang01/gsframe/loglogic"
	"testing"
)

func TestLog(t *testing.T) {
	defer loglogic.LogClose()

	defer func() {
		r := recover()
		if r != nil {
			loglogic.PError(r)

		}
	}()
	loglogic.Init(0, "logs")
	loglogic.SetListenKeyID(1001)
	loglogic.PDebug("test1")
	loglogic.PInfo("test2")
	loglogic.PInfoKey("test3", 1001)
	loglogic.PDebugKey("test4", 1002)
	loglogic.PDebugKey("test5", 1001)
	//errs := fmt.Errorf(string(debug.Stack()))
	panic("panicinfo")

}
