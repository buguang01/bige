package mysqlmodel_test

import (
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/mysqlmodel"
	"testing"
)

var (
	//不管是什么模式，都需要一个全局的变量来放连接
	DBExample *mysqlmodel.MysqlAccess
)

func init() {
	DBExample = mysqlmodel.NewMysqlAccess(&mysqlmodel.MysqlConfigModel{
		Dsn:        "root:6JkZsIybo25ls81a@tcp(192.168.39.97:3306)/test?charset=utf8",
		MaxOpenNum: 2000,
		MaxIdleNum: 1000,
	})
}

func TestDB(t *testing.T) {
	loglogic.Init(0, "logs")
	defer loglogic.LogClose()
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
