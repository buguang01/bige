package util

import (
	"runtime"
	"strings"
)

//PrintMyName 返回调用方法名
func PrintMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	s := runtime.FuncForPC(pc).Name()
	n := strings.LastIndex(s, ".")
	return s[n+1:]
}
