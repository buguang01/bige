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

//IntArrayToString []int 转string
func IntArrayToString(arr []int, sub string) string {
	sb := NewStringBuilder()
	for i, v := range arr {
		if i != 0 {
			sb.Append(sub)
		}
		sb.AppendInt(v)
	}
	return sb.ToString()
}
