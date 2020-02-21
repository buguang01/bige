package modules_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	sk    *modules.SocketModule
	skcli *modules.SocketCliModule
)

func TestService(t *testing.T) {
	Logger.Init(0, "logs", Logger.LogModeFmt)

	smd := modules.NewGameService()
	// logic = modules.NewLogicModule()
	// data = modules.NewDataBaseModule(&sql.DB{})
	// web = modules.NewWebModule()
	// wss = modules.NewWebSocketModule()
	// task = modules.NewAutoTaskModule()
	sk = modules.NewSocketModule(
		modules.SocketSetTimeout(3),
	)
	skcli = modules.NewSocketCliModule()
	// nsq = modules.NewNsqdModule()
	smd.AddModule(sk, skcli)
	smd.Run()
	Logger.LogClose()
}

func TestTask(t *testing.T) {
	Logger.Init(0, "logs", Logger.LogModeFmt)
	mod := modules.NewAutoTaskModule()
	mod.Init()
	mod.Start()
	tk := new(TaskModel)
	tk.AutoTaskModel.Handle = tk.Handle
	mod.AddTask(tk)
	time.Sleep(100 * time.Second)
}

type TaskModel struct {
	modules.AutoTaskModel
}

func (task *TaskModel) Handle(ctx context.Context) {
	fmt.Println("abc")
}
