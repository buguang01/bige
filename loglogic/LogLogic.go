package loglogic

import (
	"log"
	"os"
)

//日志等级之间都有10段预留
//为了日后实际开发时，可能出现的自定义日志等级
const (
	debuglevel  = 10
	infolevel   = 20
	statuslevel = 30
	errorlevel  = 40
	fatallevel  = 50
)

const (
	pdebuglevel  = "[debug		]"
	pinfolevel   = "[info		]"
	pstatuslevel = "[status		]"
	perrorlevel  = "[error		]"
	pfatallevel  = "[fatal		]"
)

var (
	logExample *LogLogic
)

//LogLogic 写日志用的类，里面会自行维护要写到哪个文件里去；
//第一次使用会按时间开一个当前时间的文件
//如果到了第二天，在第二天的第一次写入时，会关闭之前的文件
//重新打开一个新的文件来写
//文件会在这个基础上，再分为不同等级的日志，写在不同的文件中
//一天的日志，都会写在同一日期的文件夹下面
//如果设置了特殊监听的keyid，那会在之前的基础上，再加一个文件
//会把那个文件的名字中日志等级的部分改成这个keyid
type LogLogic struct {
	LogChan chan *LogMsgModel //写日志的信道
	Logger  *log.Logger       //写日志的系统对象
	Logfile *os.File          //对应的日志文件
}

//LogMsgModel 日志的信道数据结构
type LogMsgModel struct {
}
