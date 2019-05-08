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

func (this *SqlDataModel) GetKeyID() int {
	return this.KeyID
}
func (this *SqlDataModel) GetDataKey() string {
	return this.DataKey
}

func (this *SqlDataModel) GetUpTime() time.Duration {
	return this.UpTime
}

func (this *SqlDataModel) UpDataSave(conndb *sql.DB) error {
	if this.SaveFun != nil {
		if this.DataDBModel != nil {
			return this.SaveFun(conndb, this.DataDBModel)
		}
	}
	return nil
}

type DataDBModel interface {
}

type UpDataSave func(conndb *sql.DB, datamd DataDBModel) error

//ISqlDataModel 保存DB的接口
type ISqlDataModel interface {
	//用户主键,一般一个用户需要一个专门的协程来负责对他进行保存操作
	//这个ID就是用来确认这件事的
	GetKeyID() int
	//数据表,如果你的表放入时，不是马上保存的，那么后续可以用这个KEY来进行覆盖，
	//这样就可以实现多次修改一次保存的功能
	GetDataKey() string
	//保存时间，每次过来的时候，需要告诉我，你需要最长的保存时间，
	//如果下次保存时间比你的大，那就用你的这个时间 替换掉
	GetUpTime() time.Duration
	//这是被调用的方法
	UpDataSave(conndb *sql.DB) error
}
