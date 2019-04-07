package loglogic

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

//LogHandleModel 写日志用的类，里面会自行维护要写到哪个文件里去；
//第一次使用会按时间开一个当前时间的文件
//如果到了第二天，在第二天的第一次写入时，会关闭之前的文件
//重新打开一个新的文件来写
//文件会在这个基础上，再分为不同等级的日志，写在不同的文件中
//一天的日志，都会写在同一日期的文件夹下面
//如果设置了特殊监听的keyid，那会在之前的基础上，再加一个文件
//会把那个文件的名字中日志等级的部分改成这个keyid
type LogHandleModel struct {
	LogChan chan *LogMsgModel //写日志的信道
	Logger  *log.Logger       //写日志的系统对象
	Logfile *os.File          //对应的日志文件
	CurrDay time.Time         //写入目录的那个日期部分，用来确定是不是要新开个对象
	wg      sync.WaitGroup    //用来确认是不是关了
}

//NewLogHandle 开新的日志
func NewLogHandle(dt time.Time, lv LogLevel, pathstr string) (result *LogHandleModel) {
	result = new(LogHandleModel)

	result.LogChan = make(chan *LogMsgModel, 10)

	filename := fmt.Sprintf("%s_%02d.%02d.%02d.log",
		GetFileNameByLogLevel(lv),
		dt.Hour(),
		dt.Minute(),
		dt.Second())
	result.Logfile, _ = os.OpenFile(path.Join(pathstr, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// result.Logfile, _ = os.Create(path.Join(pathstr, filename))
	result.Logger = log.New(result.Logfile, "", 0)
	go result.handle()
	return result
}

//NewLogHandleByKeyID 开新的日志 用keyid来开
func NewLogHandleByKeyID(dt time.Time, keyid int, pathstr string) (result *LogHandleModel) {
	result = new(LogHandleModel)
	result.LogChan = make(chan *LogMsgModel, 10)

	filename := fmt.Sprintf("%d_%02d.%02d.%02d.log",
		keyid,
		dt.Hour(),
		dt.Minute(),
		dt.Second())
	result.Logfile, _ = os.OpenFile(path.Join(pathstr, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//os.Create(path.Join(pathstr, filename))
	result.Logger = log.New(result.Logfile, "", 0)
	go result.handle()
	return result
}

//handle 写日志的协程
func (lghd *LogHandleModel) handle() {
	lghd.wg.Add(1)
	defer lghd.Logfile.Close()
	defer lghd.wg.Done()
	for msg := range lghd.LogChan {
		if msg.Stack == "" {
			lghd.Logger.Output(2, fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d %s %s",
				msg.CreateTime.Year(),
				msg.CreateTime.Month(),
				msg.CreateTime.Day(),
				msg.CreateTime.Hour(),
				msg.CreateTime.Minute(),
				msg.CreateTime.Second(),
				GetLogNameByLogLevel(msg.LogLv),
				msg.Msg))
		} else {
			lghd.Logger.Println(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d %s %s\r\n%v",
				msg.CreateTime.Year(),
				msg.CreateTime.Month(),
				msg.CreateTime.Day(),
				msg.CreateTime.Hour(),
				msg.CreateTime.Minute(),
				msg.CreateTime.Second(),
				GetLogNameByLogLevel(msg.LogLv),
				msg.Msg,
				msg.Stack))
		}

	}
}

//Close 关闭本日志
func (lghd *LogHandleModel) Close() {
	close(lghd.LogChan)
}

//WaitClose 关闭并等待
func (lghd *LogHandleModel) WaitClose() {
	close(lghd.LogChan)
	lghd.wg.Wait()
}
