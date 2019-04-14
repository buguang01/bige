package util

import (
	"fmt"
	"reflect"
	"strconv"
)

//类型转换器
const Convert convert = true

type convert bool

//ToInt32 类型转到Int32
func (convert) ToInt32(v interface{}) int32 {
	switch v.(type) {
	case string:
		result, _ := strconv.ParseInt(v.(string), 10, 32)
		return int32(result)
	case float64:
		return int32(v.(float64))
	}
	panic(fmt.Sprintf("%#v to int32 fail.", v))
}

//ToString 类型转到string
func (convert) ToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.FormatInt(int64(v.(int)), 10)
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	}
	panic(fmt.Sprintf("%#v to string fail.", v))

}

func (convert) ByteToAll(src []byte, dest interface{}) {
	s := string(src)
	dpv := reflect.ValueOf(dest)
	dv := reflect.Indirect(dpv)
	switch dest.(type) {
	case *string:
		dv.SetString(s)

	case *int32, *int, *int64:
		in64, _ := strconv.ParseInt(s, 10, dv.Type().Bits())
		dv.SetInt(in64)
		// case *int:
		// 	in64, _ := strconv.ParseInt(s, 10, dv.Type().Bits())
		// 	dest = int(in64)
		// case *int64:
		// 	in64, _ := strconv.ParseInt(s, 10, dv.Type().Bits())
		// 	dest = int64(in64)
	}
}
