package util

import (
	"strings"
)

//StringToIntArray 字符串转int数组
func StringToIntArray(str string, sub string) []int {
	arr := strings.Split(str, sub)
	result := make([]int, len(arr))
	for i, v := range arr {
		result[i], _ = NewString(v).ToInt()
	}
	return result
}
