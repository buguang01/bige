package util

import (
	"strings"
)

//StringToInt32Array 字符串转int32数组
func StringToInt32Array(str string, sub string) []int32 {
	arr := strings.Split(str, sub)
	result := make([]int32, len(arr))
	for i, v := range arr {
		result[i] = Convert.ToInt32(v)
	}
	return result
}
