package model

import (
	"github.com/buguang01/Logger"
	"database/sql"
	"time"

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
	Dsn         string        //数据库连接字符串
	MaxOpenNum  int           //最大连接数
	MaxIdleNum  int           //最大空闲连接数
	MaxLifetime time.Duration //连接空闲等待时间（秒）
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
	result.DBConobj.SetMaxOpenConns(cgmodel.MaxOpenNum)
	result.DBConobj.SetMaxIdleConns(cgmodel.MaxIdleNum)
	result.DBConobj.SetConnMaxLifetime(cgmodel.MaxLifetime * time.Second)
	if err != nil {
		Logger.PFatal(err)
		panic(err)
	}
	Logger.PDebug("mysql init.")
	return result
}

//Ping 确认一下数据库连接
func (access *MysqlAccess) Ping() error {
	return access.DBConobj.Ping()
}

//GetConnBegin 拿到事件连接对象，不用的时候需要执行 Commit()或Rollback()
func (access *MysqlAccess) GetConnBegin() *sql.Tx {
	result, err := access.DBConobj.Begin()
	if err != nil {
		// Logger.PFatal(err)
		panic(err)
	}
	return result
}

//GetDB 拿到的并不是具体的连接，但你使用的时候，他会去池子里找个连接给你
func (access *MysqlAccess) GetDB() *sql.DB {
	return access.DBConobj
}

//Close 关闭池子,只有关服的时候，才会用到这个，一般不用也没有关系，也会自己关闭的
func (access *MysqlAccess) Close() {
	access.DBConobj.Close()
	Logger.PDebug("mysql close.")
}

//NewRead 出一个读取器
func NewRead(rows *sql.Rows) *ReadRow {
	result := new(ReadRow)
	result.Rows = rows
	result.Columns, _ = rows.Columns()
	result.Data = make([]interface{}, len(result.Columns))
	return result
}

//ReadRow 行读取器
type ReadRow struct {
	Rows    *sql.Rows
	Columns []string      //列
	Data    []interface{} //数据

}

//Read 读下一行
func (read *ReadRow) Read() bool {
	ok := read.Rows.Next()
	scanargs := make([]interface{}, len(read.Columns))
	for i := 0; i < len(read.Columns); i++ {
		scanargs[i] = &read.Data[i]
	}
	if ok {
		read.Rows.Scan(scanargs...)
	}
	return ok
}

//GetRowByColName 拿当前行的指定列
func (read *ReadRow) GetRowByColName(colname string) []byte {
	for i, col := range read.Columns {
		if col == colname {
			return read.Data[i].([]byte)
		}
	}
	return nil
}
