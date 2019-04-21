package event

import (
	"database/sql"
	"time"
)

type SqlDataModel struct {
	KeyID       int           //用户主键
	DataKey     string        //数据表
	UpTime      time.Duration //保存时间
	SaveFun     UpDataSave    //保存方法
	DataDBModel DataDBModel   //要保存的东西
}

type DataDBModel interface {
}

type UpDataSave func(conndb *sql.DB, datamd DataDBModel) error
