package modules

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util"
	"github.com/buguang01/util/threads"
)

type DataBaseModule struct {
	conndb    *sql.DB                          //数据库连接对象
	chanNum   int                              //通道缓存空间
	timeout   time.Duration                    //超时时间
	logicList map[int]*dataBaseThread          //子逻辑列表
	keyList   []int                            //key列表，用来间隔遍历
	chanList  chan []messages.IDataBaseMessage //消息信道
	getNum    int64                            //收到的总消息数
	saveNum   int64                            //保存次数
	thgo      *threads.ThreadGo                //子协程管理器
}

func NewDataBaseModule(conndb *sql.DB, opts ...options) *DataBaseModule {
	result := &DataBaseModule{
		conndb:    conndb,
		chanNum:   1024,
		timeout:   2 * time.Minute,
		logicList: make(map[int]*dataBaseThread, moduleCap),
		keyList:   make([]int, 0, moduleCap),
		getNum:    0,
		saveNum:   0,
		thgo:      threads.NewThreadGo(),
	}
	return result
}

//Init 初始化
func (mod *DataBaseModule) Init() {
	mod.chanList = make(chan []messages.IDataBaseMessage, mod.chanNum)
}

//Start 启动
func (mod *DataBaseModule) Start() {
	mod.thgo.Go(mod.Handle)
	Logger.PStatus("DataBase Module Start!")
}

//Stop 停止
func (mod *DataBaseModule) Stop() {
	/*
		当这个服务可以被停止的时候，外部就不再会有发消息进来了
		所以这就可以直接关闭通道
	*/
	close(mod.chanList)
	mod.thgo.CloseWait()
	Logger.PStatus("DataBase Module Shop!")
}

//PrintStatus 打印状态
func (mod *DataBaseModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tDataBase Module\t:%d/%d/%d\t(logic/get/save)",
		len(mod.logicList),
		atomic.LoadInt64(&mod.getNum),
		atomic.LoadInt64(&mod.saveNum))
}

func (mod *DataBaseModule) Handle(ctx context.Context) {
	tk := time.NewTicker(1 * time.Second)
	defer tk.Stop()
	loop := 0
	for {
		select {
		case msgs, ok := <-mod.chanList:
			{
				if !ok {
					//通道如果被关闭了，就可以关闭子协程了
					for _, lth := range mod.logicList {
						lth.stop()
					}
					return
				}
				if len(msgs) == 0 {
					continue
				}
				upmd := msgs[0]

				lth, ok := mod.logicList[upmd.DBThreadID()]
				if !ok {
					//新开一个协程
					lth = newDataBaseThread(upmd.DBThreadID(), mod.chanNum, mod.conndb)
					mod.logicList[lth.DBThreadID] = lth
					mod.keyList = append(mod.keyList, lth.DBThreadID)
					lth.start(mod)
				}
				lth.addMsg(msgs)
			}
		case <-tk.C:
			{
				lilen := len(mod.keyList)
				if lilen == 0 {
					continue
				}
				loop = loop % lilen
				keyid := mod.keyList[loop]
				if lth, ok := mod.logicList[keyid]; ok {
					if lth.GetMsgNum() == 0 &&
						util.GetCurrTime().Sub(lth.upTime) > mod.timeout {
						//确定子协程可以关闭
						lth.stop()
						delete(mod.logicList, keyid)
						mod.keyList = append(mod.keyList[:loop], mod.keyList[loop+1:]...)
					}
				}
				loop++
			}
		}
	}
}

func (mod *DataBaseModule) AddMsg(msgs ...messages.IDataBaseMessage) {
	atomic.AddInt64(&mod.getNum, 1)
	mod.chanList <- msgs
}

type dataBaseThread struct {
	DBThreadID int                                  //协程ID
	upDataList map[string]messages.IDataBaseMessage //缓存要更新的数据
	chanList   chan []messages.IDataBaseMessage     //收要更新的数据
	Conndb     *sql.DB                              //数据库连接对象
	upTime     time.Time                            //更新时间
	cancel     context.CancelFunc                   //关闭
}

func newDataBaseThread(id, channum int, conn *sql.DB) *dataBaseThread {
	result := &dataBaseThread{
		DBThreadID: id,
		upDataList: make(map[string]messages.IDataBaseMessage),
		chanList:   make(chan []messages.IDataBaseMessage, channum),
		Conndb:     conn,
	}
	return result
}

func (lth *dataBaseThread) start(mod *DataBaseModule) {
	lth.cancel = mod.thgo.SubGo(
		func(ctx context.Context) {
			lth.handle(ctx, mod)
		},
	)
}

func (lth *dataBaseThread) stop() {
	lth.cancel()
	close(lth.chanList)
}

func (lth *dataBaseThread) handle(ctx context.Context, mod *DataBaseModule) {
	tk := time.NewTimer(time.Second)
	defer tk.Stop()
	isruned := false
trheadhandle:
	for {
		select {
		case msg, ok := <-lth.chanList:
			{
				if !ok {
					lth.save()
					atomic.AddInt64(&mod.saveNum, 1)
					break trheadhandle
				}
				if len(msg) == 0 {
					continue
				}
				for _, data := range msg {
					lth.upDataList[data.GetDataKey()] = data
				}
				if isruned {
					tk.Reset(time.Second)
					isruned = false
				}

			}
		case <-tk.C:
			{
				if len(lth.upDataList) > 0 {
					lth.save()
					atomic.AddInt64(&mod.saveNum, 1)
					lth.upDataList = make(map[string]messages.IDataBaseMessage)
				}
				isruned = true
			}
		}
	}
}

func (lth *dataBaseThread) addMsg(msgs []messages.IDataBaseMessage) {
	lth.upTime = util.GetCurrTime()
	lth.chanList <- msgs
}

func (lth *dataBaseThread) save() {
	if tx, err := lth.Conndb.Begin(); err == nil {
		threads.Try(func() {
			for _, data := range lth.upDataList {
				if err = data.SaveDB(tx); err != nil {
					panic(errors.New(fmt.Sprintf(" keyid:%d;DataKey:%s; ", data.DBThreadID(), data.GetDataKey())))
				}
			}
			tx.Commit()
		}, func(err interface{}) {
			tx.Rollback()
			Logger.PFatal(err)
		}, nil)
	}
}

//还有多少消息没有处理完
func (lth *dataBaseThread) GetMsgNum() int {
	return len(lth.chanList) + len(lth.upDataList)
}
