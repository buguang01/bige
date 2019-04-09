package mysqlmodel

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//MysqlAccess mysql连接器
type MysqlAccess struct {
	DBConobj *sql.DB //数据库的连接池对象

}
