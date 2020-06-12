package model_test

import (
	"fmt"
	"testing"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/model"
)

func TestRedis(t *testing.T) {
	Logger.Init(0, "", Logger.LogModeFmt)
	defer Logger.LogClose()
	rd := model.NewRedisAccess()
	rdmd := rd.GetConn()
	// rdmd.Scan(0, "", 10)
	if result, err := rdmd.SrandMember("fruit", 2); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}
