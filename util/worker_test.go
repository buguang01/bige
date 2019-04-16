package util_test

import (
	"buguang01/gsframe/util"
	"fmt"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	g := util.NewIDGenerator()
	g.SetWorkerId(1001)
	g.Init()
	for i := 0; i < 100; i++ {
		go func() {
			for {
				fmt.Print(g.NextId())
				fmt.Print(",")
			}
		}()
	}
	time.Sleep(100 * time.Second)
}
