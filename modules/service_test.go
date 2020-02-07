package modules_test

import (
	"database/sql"
	"testing"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/modules"
)

var (
	web   *modules.WebModule
	wss   *modules.WebSocketModule
	nsq   *modules.NsqdModule
	logic *modules.LogicModule
	data  *modules.DataBaseModule
	task  *modules.AutoTaskModule
)

func TestService(t *testing.T) {
	Logger.Init(0, "logs", Logger.LogModeFmt)

	smd := modules.NewGameService()
	logic = modules.NewLogicModule()
	data = modules.NewDataBaseModule(&sql.DB{})
	web = modules.NewWebModule()
	wss = modules.NewWebSocketModule()
	task = modules.NewAutoTaskModule()
	// nsq = modules.NewNsqdModule()
	smd.AddModule(data, logic, task, web, wss)

	smd.Run()
	Logger.LogClose()
}
