package loglogic

import (
	"time"
)

//LogMsgModel 日志的信道数据结构
type LogMsgModel struct {
	Msg        string    //日志文本
	LogLv      LogLevel  //日志等级
	KeyID      int       //可能有的日志写入KEYID，-1为不用管
	Stack      string    //堆栈信息
	CreateTime time.Time //日志生成的时间
}
