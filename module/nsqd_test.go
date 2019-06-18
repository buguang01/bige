package module_test

import (
	"testing"
	"time"

	"github.com/buguang01/Logger"

	"github.com/buguang01/gsframe/module"
)

func TestNsqd(t *testing.T) {
	Logger.Init(1, "logs")

	nsqd := module.NewNsqdModule(&module.NsqdConfig{
		Addr:                []string{"192.168.39.97:4150"},
		NSQLookupdAddr:      []string{"192.168.39.97:4161"},
		ChanNum:             100,
		LookupdPollInterval: 1000,
		MaxInFlight:         1000,
	}, 4)
	nsqd.Init()
	nsqd.Start()

	time.Sleep(100 * time.Second)
}
