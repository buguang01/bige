package model_test

import (
	"fmt"
	"testing"

	"github.com/buguang01/Logger"
	"github.com/buguang01/gsframe/model"
)

var (
	//不管是什么模式，都需要一个全局的变量来放连接
	DBExample *model.MysqlAccess
)

func init() {
	// DBExample = model.NewMysqlAccess(&model.MysqlConfigModel{
	// 	Dsn:        "root:6JkZsIybo25ls81a@tcp(192.168.39.97:3306)/test?charset=utf8",
	// 	MaxOpenNum: 2000,
	// 	MaxIdleNum: 1000,
	// })
}

//写事务的方式
func TestDBTran(t *testing.T) {
	Logger.Init(0, "logs", Logger.LogModeFmt)
	defer Logger.LogClose()
	db := DBExample.GetConnBegin()
	defer func() {
		if err := recover(); err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	r, err := db.Exec("insert into abtable (name) values(?)", "xiacs5")
	_ = r
	_ = err
	db.Commit()
	r, err = db.Exec("insert into abtable (name) values(?)", "xiacs6")
}

func TestQuery(t *testing.T) {
	Logger.Init(Logger.LogLeveldebuglevel, "", Logger.LogModeFmt)
	defer Logger.LogClose()
	cf := model.RedisConfigModel{
		ConAddr:  "152.136.222.222:6379",
		Password: "cMz8eEv3fT0XD2ue",
	}
	redis := model.NewRedisAccess(&cf)
	rd := redis.GetConn()
	arr, _ := rd.DictGetAllByStringArray("testlist")
	fmt.Println(arr)

}
