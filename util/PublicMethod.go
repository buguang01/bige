package util

import "runtime"

//PrintMyName 返回调用方法名
func PrintMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
