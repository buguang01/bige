package modules_test

import (
	"testing"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/modules"
)

var (
	web *modules.WebModule
	wss *modules.WebSocketModule
	nsq *modules.NsqdModule
)

func TestService(t *testing.T) {
	Logger.Init(0, "logs", Logger.LogModeFmt)

	smd := modules.NewGameService()
	web = modules.NewWebModule()
	wss = modules.NewWebSocketModule()
	nsq = modules.NewNsqdModule()
	smd.AddModule(web, wss, nsq)
	smd.Run()

}
