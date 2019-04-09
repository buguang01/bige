package mysqlmodel

import (
	"buguang01/gsframe/loglogic"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //注册MYSQL
)

//MysqlConfigModel 数据库地配置
//下面是4种连接字符串的写法
// user@unix(/path/to/socket)/dbname?charset=utf8
// user:password@tcp(localhost:5555)/dbname?charset=utf8
// user:password@/dbname
// user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname
//
// 如果是给http的服务用，建议连连数大一些，空闲也要大
// 如果是给专门的DB模块使用，这二个数就按你自己对应这个模块的协程数来定；
// 差不多1：1就可以了。
type MysqlConfigModel struct {
	Dsn        string //数据库连接字符串
	MaxOpenNum int32  //最大连接数
	MaxIdleNum int32  //最大空闲连接数
}

//MysqlAccess mysql连接器
type MysqlAccess struct {
	DBConobj *sql.DB           //数据库的连接池对象
	cg       *MysqlConfigModel //配置信息
}

//NewMysqlAccess 新建一个数据库管理器
func NewMysqlAccess(cgmodel *MysqlConfigModel) *MysqlAccess {
	var err error
	result := new(MysqlAccess)
	result.cg = cgmodel
	result.DBConobj, err = sql.Open("mysql", cgmodel.Dsn)
	if err != nil {
		loglogic.PFatal(err)
		panic(err)
	}
	return result
}

//GetConnBegin 拿到事件连接对象，不用的时候需要执行 Commit()或Rollback()
func (access *MysqlAccess) GetConnBegin() *sql.Tx {
	result, err := access.DBConobj.Begin()
	if err != nil {
		// loglogic.PFatal(err)
		panic(err)
	}
	return result
}
