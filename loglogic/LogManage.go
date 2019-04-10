package loglogic

//有一个日志的协程用来收消息，然后再发给具体的写文件的协程
import (
	"buguang01/gsframe/util"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"sync"
	"time"
)

var (
	logExample *LogManageModel
)

// func init() {
// 	// logExample = &LogManageModel{
// 	// 	MiniLv:        LogLeveldebuglevel,
// 	// 	PathStr:       "logs",
// 	// 	LogList:       make(map[LogLevel]*LogHandleModel),
// 	// 	LogKeyList:    make(map[int]*LogHandleModel),
// 	// 	ListenKeyList: make(map[int]bool),
// 	// 	msgchan:       make(chan *LogMsgModel, 10),
// 	// }
// }

//LogManageModel 日志管理器
type LogManageModel struct {
	LogList       map[LogLevel]*LogHandleModel //当前打开的日志文件
	LogKeyList    map[int]*LogHandleModel      //当前为所有需要记录keyid日志
	MiniLv        LogLevel                     //最小等级
	PathStr       string                       //保存的基础目录
	CurrDir       string                       //是按时间做目录
	CurrDay       time.Time                    //当前写日志的那个目录的日期
	ListenKeyList map[int]bool                 //监听key列表
	msgchan       chan *LogMsgModel            // 日志扇入流
	wg            sync.WaitGroup               //用来确认是不是关了

}

//Init 初始化日志管理器的参数
func Init(minlv LogLevel, pathstr string) {
	logExample = new(LogManageModel)
	logExample.MiniLv = minlv
	logExample.PathStr = pathstr
	logExample.CurrDay = util.TimeConvert.GetMinDateTime()
	logExample.LogList = make(map[LogLevel]*LogHandleModel)
	logExample.LogKeyList = make(map[int]*LogHandleModel)
	logExample.ListenKeyList = make(map[int]bool)
	logExample.msgchan = make(chan *LogMsgModel, 10)
	go logExample.Handle()
}

//LogClose 正常关闭日志服务； 在程序退出的时候，才可以运行的方法
func LogClose() {
	close(logExample.msgchan)
	logExample.wg.Wait()
	for _, lghd := range logExample.LogKeyList {
		lghd.WaitClose()
	}
	for _, lghd := range logExample.LogList {
		lghd.WaitClose()
	}
	// time.Sleep(10 * time.Second)
}

//Handle logmanage的主协程
func (lgmd *LogManageModel) Handle() {
	lgmd.wg.Add(1)
	defer lgmd.wg.Done()
	for msg := range lgmd.msgchan {
		//拿到了一个要写的日志
		var lghd *LogHandleModel
		//确认目录是不是有了
		lgmd.checkDir(msg.CreateTime)
		//如果有设置keyid就检查一下这个keyid要不要写日志
		if msg.KeyID != -1 {
			if _, ok := lgmd.ListenKeyList[msg.KeyID]; ok {
				lghd = lgmd.getLogHandleByKeyID(msg.KeyID)
				lghd.LogChan <- msg
			}
		}
		//判断是不是需要写入的日志
		if msg.LogLv < lgmd.MiniLv {
			continue
		}

		//拿到目标日志对象
		lghd = lgmd.getLogModel(msg.LogLv)
		lghd.LogChan <- msg
		//写入主日志文件
		lghd = lgmd.getLogModel(0)
		lghd.LogChan <- msg

	}
}

//checkDir  检查只定新时间是否需要修改目录
func (lgmd *LogManageModel) checkDir(d time.Time) {
	if lgmd.CurrDay == util.TimeConvert.GetDate(d) {
		return
	}
	lgmd.CurrDay = util.TimeConvert.GetCurrDate()
	dir := fmt.Sprintf("%d_%02d_%02d",
		lgmd.CurrDay.Year(),
		lgmd.CurrDay.Month(),
		lgmd.CurrDay.Day())
	lgmd.CurrDir = path.Join(lgmd.PathStr, dir)
	os.MkdirAll(lgmd.CurrDir, os.ModePerm)

}

//getLogModel 拿到要写日志的对象
func (lgmd *LogManageModel) getLogModel(lv LogLevel) (result *LogHandleModel) {

	result, ok := lgmd.LogList[lv]
	if !ok {
		//基本就是拿不到那个对象，需要新建一个
		result = NewLogHandle(util.TimeConvert.GetCurrTime(), lv, lgmd.CurrDir)
		lgmd.LogList[lv] = result
	} else if result.CurrDay != lgmd.CurrDay {
		result.Close()
		result = NewLogHandle(util.TimeConvert.GetCurrTime(), lv, lgmd.CurrDir)
		lgmd.LogList[lv] = result
	}
	return result
}

//getLogHandleByKeyID 拿到要写日志的对象，指定KEYID
func (lgmd *LogManageModel) getLogHandleByKeyID(keyid int) (result *LogHandleModel) {
	result, ok := lgmd.LogKeyList[keyid]
	if !ok {
		//基本就是拿不到那个对象，需要新建一个
		result = NewLogHandleByKeyID(util.TimeConvert.GetCurrTime(), keyid, lgmd.CurrDir)
		lgmd.LogKeyList[keyid] = result
	} else if result.CurrDay != lgmd.CurrDay {
		result.Close()
		result = NewLogHandleByKeyID(util.TimeConvert.GetCurrTime(), keyid, lgmd.CurrDir)
		lgmd.LogKeyList[keyid] = result
	}
	return result
}

//SetListenKeyID 设置监听keyid
func SetListenKeyID(keyid int) {
	logExample.ListenKeyList[keyid] = true
}

//RemoveListenKeyID 移除监听的keyid
func RemoveListenKeyID(keyid int) {
	delete(logExample.ListenKeyList, keyid)
	lghd, ok := logExample.LogKeyList[keyid]
	if ok {
		lghd.Close()
	}
}

//PDebug 调试日志
func PDebug(msgstr string, a ...interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf(msgstr, a...),
		LogLv: LogLeveldebuglevel,
		Stack: "",
		KeyID: -1,
	})
}

//PDebugKey 指定key的调试日志
func PDebugKey(msgstr string, keyid int, a ...interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf(msgstr, a...),
		LogLv: LogLeveldebuglevel,
		Stack: "",
		KeyID: keyid,
	})
}

//PInfo 一般日志
func PInfo(msgstr string, a ...interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf(msgstr, a...),
		LogLv: LogLevelinfolevel,
		Stack: "",
		KeyID: -1,
	})
}

//PInfoKey 指定key的一般日志
func PInfoKey(msgstr string, keyid int, a ...interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf(msgstr, a...),
		LogLv: LogLevelinfolevel,
		Stack: "",
		KeyID: keyid,
	})
}

//PStatus 服务器状态日志
func PStatus(msgstr string, a ...interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf(msgstr, a...),
		LogLv: LogLevelstatuslevel,
		Stack: "",
		KeyID: -1,
	})
}

//PStatusKey 指定key的服务器状态日志
func PStatusKey(msgstr string, keyid int, a ...interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf(msgstr, a...),
		LogLv: LogLevelstatuslevel,
		Stack: "",
		KeyID: keyid,
	})
}

//PError 游戏中的错误日志
func PError(msgstr interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf("%v", msgstr),
		LogLv: LogLevelerrorlevel,
		Stack: string(debug.Stack()),
		KeyID: -1,
	})
}

//PErrorKey 指定key的游戏中的错误日志
func PErrorKey(msgstr interface{}, keyid int) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf("%v", msgstr),
		LogLv: LogLevelerrorlevel,
		Stack: string(debug.Stack()),
		KeyID: keyid,
	})
}

//PFatal 程序异常日志
func PFatal(msgstr interface{}) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf("%v", msgstr),
		LogLv: LogLevelfatallevel,
		Stack: string(debug.Stack()),
		KeyID: -1,
	})
}

//PFatalKey 指定key的程序异常日志
func PFatalKey(msgstr interface{}, keyid int) {
	PrintLog(&LogMsgModel{
		Msg:   fmt.Sprintf("%v", msgstr),
		LogLv: LogLevelfatallevel,
		Stack: string(debug.Stack()),
		KeyID: keyid,
	})
}

//PrintLog 扇入的入口
func PrintLog(msg *LogMsgModel) {
	msg.CreateTime = util.TimeConvert.GetCurrTime()
	logExample.msgchan <- msg
}
