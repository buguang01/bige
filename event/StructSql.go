package event

import (
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/buguang01/bige/model"
	"github.com/buguang01/util"
)

//按接口查询数据
func DataGetByID(conndb *sql.DB, v IStructSql) *sql.Rows {
	sqlstr := MarshalQsql(v, v.GetTableName())
	result, err := conndb.Query(sqlstr, v.QueryArray())
	if err != nil {
		panic(err)
	}
	return result
}

//生成更新SQL
func MarshalUpSql(v interface{}, tablename string) (sql string) {
	result := util.NewStringBuilder()
	result.Append("INSERT INTO ")
	result.Append(tablename)
	result.Append("(")
	t := reflect.TypeOf(v)
	farr := t.Elem()
	tmp := util.NewStringBuilder()
	vtmp := util.NewStringBuilder()
Fieldfor:
	for i := 0; i < farr.NumField(); i++ {
		field := farr.Field(i)
		bigetag := field.Tag.Get("bige")
		narr := strings.Split(bigetag, ",")
		name := field.Name
		iskey := false
		for _, v := range narr {
			switch v {
			case "bigekey":
				iskey = true
			case "select":
			case "-":
				continue Fieldfor
			default:
				name = v
			}
		}
		if !tmp.IsEmpty() {
			result.Append(",")
			tmp.Append(",")
		}
		result.Append(name)
		tmp.Append("?")
		if !iskey {
			if !vtmp.IsEmpty() {
				vtmp.Append(",")
			}
			vtmp.Append(name)
			vtmp.Append("=values(")
			vtmp.Append(name)
			vtmp.Append(")")
		}
	}
	result.Append(")VALUES(")
	result.Append(tmp.ToString())
	result.Append(") ON DUPLICATE KEY UPDATE ")
	result.Append(vtmp.ToString())
	result.Append(";")
	return result.ToString()
}

//生成查询SQL
func MarshalQsql(v interface{}, tablename string) (sql string) {
	result := util.NewStringBuilder()
	result.Append("SELECT ")

	t := reflect.TypeOf(v)
	farr := t.Elem()
	where := util.NewStringBuilder()
Fieldfor:
	for i := 0; i < farr.NumField(); i++ {
		if i > 0 {
			result.Append(",")
		}
		field := farr.Field(i)
		bigetag := field.Tag.Get("bige")
		narr := strings.Split(bigetag, ",")
		name := field.Name
		iswhere := false
		for _, v := range narr {
			switch v {
			case "bigekey":
			case "select":
				iswhere = true
			case "-":
				continue Fieldfor
			default:
				name = v
			}
		}
		result.Append(name)
		if iswhere {
			if !where.IsEmpty() {
				where.Append(" AND ")
			}
			where.Append(name)
			where.Append("=?")
		}
	}
	result.Append(" FROM ")
	result.Append(tablename)
	result.Append(" WHERE ")
	result.Append(where.ToString())
	result.Append(";")
	return result.ToString()
}

//数据映射接口
type IStructSql interface {
	//表名
	GetTableName() string
	//保存参数列表
	ParmArray() []interface{}
	//查询的参数列表
	QueryArray() []interface{}
}

//数据映射接口的导入结构
type SqlDataStructModel struct {
	KeyID       int           //用户主键
	DataKey     string        //数据表
	UpTime      time.Duration //保存时间
	DataDBModel IStructSql    //要保存的东西
}

func (this *SqlDataStructModel) GetKeyID() int {
	return this.KeyID
}
func (this *SqlDataStructModel) GetDataKey() string {
	return this.DataKey
}

func (this *SqlDataStructModel) GetUpTime() time.Duration {
	return this.UpTime
}

func (this *SqlDataStructModel) UpDataSave(conndb model.IConnDB) error {
	if this.DataDBModel != nil {
		sqlstr := MarshalUpSql(this.DataDBModel, this.DataDBModel.GetTableName())
		_, err := conndb.Exec(sqlstr, this.DataDBModel.ParmArray())
		return err
	}
	return nil
}
