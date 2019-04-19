package loglogic

import (
	"fmt"
)

//LogLevel 日志等级
type LogLevel uint

//日志等级之间都有10段预留
//为了日后实际开发时，可能出现的自定义日志等级
const (
	LogLevelmainlevel   LogLevel = 0
	LogLeveldebuglevel  LogLevel = 10
	LogLevelinfolevel   LogLevel = 20
	LogLevelstatuslevel LogLevel = 30
	LogLevelerrorlevel  LogLevel = 40
	LogLevelfatallevel  LogLevel = 50
)

const (
	pmainlevel   = "[main		]"
	pdebuglevel  = "[debug		]"
	pinfolevel   = "[info		]"
	pstatuslevel = "[status		]"
	perrorlevel  = "[error		]"
	pfatallevel  = "[fatal		]"
)

//GetFileNameByLogLevel 按日志等级生成前缀
func GetFileNameByLogLevel(lv LogLevel) (result string) {
	if lv == LogLevelmainlevel {
		result = fmt.Sprintf("main_%d", lv)
	} else if lv <= LogLeveldebuglevel {
		result = fmt.Sprintf("debug_%d", lv)
	} else if lv <= LogLevelinfolevel {
		result = fmt.Sprintf("info_%d", lv)
	} else if lv <= LogLevelstatuslevel {
		result = fmt.Sprintf("status_%d", lv)
	} else if lv <= LogLevelerrorlevel {
		result = fmt.Sprintf("error_%d", lv)
	} else if lv <= LogLevelfatallevel {
		result = fmt.Sprintf("fatal_%d", lv)
	}
	return result
}

//GetLogNameByLogLevel 按日志等级生成日志等级的名字
func GetLogNameByLogLevel(lv LogLevel) (result string) {
	if lv == LogLevelmainlevel {
		result = fmt.Sprintf("[main_%d]", lv)
	} else if lv <= LogLeveldebuglevel {
		result = fmt.Sprintf("[debug_%d]", lv)
	} else if lv <= LogLevelinfolevel {
		result = fmt.Sprintf("[info_%d]", lv)
	} else if lv <= LogLevelstatuslevel {
		result = fmt.Sprintf("[status_%d]", lv)
	} else if lv <= LogLevelerrorlevel {
		result = fmt.Sprintf("[error_%d]", lv)
	} else if lv <= LogLevelfatallevel {
		result = fmt.Sprintf("[fatal_%d]", lv)
	}
	return result
}
