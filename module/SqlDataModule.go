package module

import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/threads"
	"buguang01/gsframe/util"
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"
)

type SqlDataConfig struct {
	Timeout    int //超时时间（秒）
	InitNum    int //初始化内存空间
	ChanNum    int //通道缓存空间
	SubChanNum int //子通道缓存空间
}

//SqlDataModule 数据库的模块
type SqlDataModule struct {
	playerlist map[int]*DataThread      //操作数据库的列表
	keylist    []int                    //key列表，用来间隔遍历
	chandata   chan event.ISqlDataModel //消息信道
	mgGo       *threads.ThreadGo        //子协程管理器
	getnum     int64                    //收到的消息数
	conndb     *sql.DB                  //数据库对象
	cg         *SqlDataConfig
}

func NewSqlDataModule(config *SqlDataConfig, sqldb *sql.DB) *SqlDataModule {
	result := new(SqlDataModule)
	result.cg = config
	result.playerlist = make(map[int]*DataThread, config.InitNum)
	result.keylist = make([]int, 0, config.InitNum)
	result.chandata = make(chan event.ISqlDataModel, config.ChanNum)
	result.mgGo = threads.NewThreadGo()
	result.conndb = sqldb
	return result
}

func (this *SqlDataModule) Init() {

}
func (this *SqlDataModule) Start() {
	this.mgGo.Go(this.Handle)
	loglogic.PStatus("SqlData Module Start!")

}
func (this *SqlDataModule) Stop() {
	this.mgGo.CloseWait()
	loglogic.PStatus("SqlData Module Start!")
}

//PrintStatus 打印状态
func (this *SqlDataModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n		Data Module:	%d/%d	(RunNum/getnum)",
		len(this.playerlist),
		atomic.AddInt64(&this.getnum, 0))
}
func (this *SqlDataModule) AddMsg(msg event.ISqlDataModel) {
	// atomic.AddInt64(&this.currnum, 1)
	this.chandata <- msg
}

func (this *SqlDataModule) Handle(ctx context.Context) {
	tk := time.NewTimer(1 * time.Second)
	loop := 0
	for {
		select {
		case <-ctx.Done():
			{
				for {
					select {
					case upmd := <-this.chandata:
						{
							logicth, ok := this.playerlist[upmd.GetKeyID()]
							if !ok {
								//新开一个协程
								logicth = newDataThread(upmd.GetKeyID(), this.conndb, this.cg.SubChanNum)
								this.playerlist[logicth.KeyID] = logicth
								logicth.Start(this)
							}
							logicth.UpTime = util.GetCurrTime()
							logicth.SendChan <- upmd
						}
					default:
						{
							close(this.chandata)
							for _, logicth := range this.playerlist {
								logicth.CloseChan()
							}
							return
						}
					}
				}
			}
		case upmd, ok := <-this.chandata:
			{
				if !ok {
					return
				}
				atomic.AddInt64(&this.getnum, 1)
				logicth, ok := this.playerlist[upmd.GetKeyID()]
				if !ok {
					//新开一个协程
					logicth = newDataThread(upmd.GetKeyID(), this.conndb, this.cg.SubChanNum)
					this.playerlist[logicth.KeyID] = logicth
					this.keylist = append(this.keylist, logicth.KeyID)
					logicth.Start(this)
				}
				logicth.UpTime = util.GetCurrTime()
				logicth.SendChan <- upmd
			}
		case <-tk.C:
			{
				tk.Reset(1 * time.Second)
				if len(this.keylist) == 0 {
					break
				}
				loop := (loop + 1) % len(this.keylist)
				keyid := this.keylist[loop]
				logicth, ok := this.playerlist[keyid]
				if ok {
					if logicth.UpTime.Add(time.Duration(this.cg.Timeout)*time.Second).Unix() < util.GetCurrTime().Unix() {
						logicth.Stop(this, loop)
						loop--
					}
				} else {
					this.keylist = append(this.keylist[:loop], this.keylist[loop+1:]...)
				}
			}
		}
	}
}

//DataThread 用户的数据库协程
type DataThread struct {
	KeyID     int                            //用户主键
	updatamap map[string]event.ISqlDataModel //缓存要更新的数据
	SendChan  chan event.ISqlDataModel       //收要更新的数据
	Conndb    *sql.DB                        //数据库连接对象
	UpTime    time.Time                      //更新时间
}

func newDataThread(keyid int, conndb *sql.DB, channum int) *DataThread {
	result := new(DataThread)
	result.KeyID = keyid
	result.updatamap = make(map[string]event.ISqlDataModel)
	result.SendChan = make(chan event.ISqlDataModel, channum)
	result.Conndb = conndb

	return result
}

func (this *DataThread) Start(mg *SqlDataModule) {
	mg.mgGo.Go(
		func(ctx context.Context) {
			this.Handle(ctx, mg)
		})
}

func (this *DataThread) Stop(mg *SqlDataModule, index int) {
	delete(mg.playerlist, this.KeyID)
	mg.keylist = append(mg.keylist[:index], mg.keylist[index+1:]...)
	close(this.SendChan)
}
func (this *DataThread) CloseChan() {
	close(this.SendChan)
}

//Handle 保存数据的协程
func (this *DataThread) Handle(ctx context.Context, mg *SqlDataModule) {
	dt := time.Now()
	tk := time.NewTimer(time.Second * 600)
	tk.Stop()
threadhandle:
	for {
		select {
		case <-ctx.Done():
			{
				for upmd := range this.SendChan {
					this.updatamap[upmd.GetDataKey()] = upmd
				}
				threads.Try(this.Save, nil, nil)
				break threadhandle
			}
		case upmd, ok := <-this.SendChan:
			{
				if !ok {
					threads.Try(this.Save, nil, nil)
					break threadhandle
				}
				this.updatamap[upmd.GetDataKey()] = upmd
				nd := upmd.GetUpTime()
				if nd <= 0 {
					threads.Try(this.Save, nil, nil)
					// atomic.AddInt64(&mod.savenum, 1)
					tk.Stop()
				} else {
					ndt := time.Now().Add(nd)
					if dt.Unix() < time.Now().Unix() {
						//比当前时间小，说明可能是停下来了要重新设置时间
						dt = ndt
						tk.Reset(nd)
					} else if dt.Unix() > ndt.Unix() {
						dt = ndt
						tk.Reset(nd)
					}
				}
			}
		case <-tk.C:
			{
				threads.Try(this.Save, nil, nil)
				tk.Stop()
			}
		}
	}

	// atomic.AddInt64(&mod.currnum, 1)
	// d := time.Second * 600
	// dt := time.Now().Add(d)
	// tk := time.NewTimer(d)
	// tk.Stop()
	// runing := true
	// for runing {
	// 	threads.Try(func() {
	// 		select {
	// 		case <-this.ctx.Done():
	// 			for upmd := range this.SendChan {
	// 				this.updatamap[upmd.DataKey] = upmd
	// 			}
	// 			this.Save()
	// 			atomic.AddInt64(&mod.savenum, 1)
	// 			runing = false
	// 			return
	// 		case <-tk.C:
	// 			this.Save()
	// 			atomic.AddInt64(&mod.savenum, 1)
	// 			tk.Stop()

	// 		case upmd, ok := <-this.SendChan:
	// 			if !ok {
	// 				this.Save()
	// 				atomic.AddInt64(&mod.savenum, 1)
	// 				runing = false
	// 				return
	// 			}
	// 			this.updatamap[upmd.DataKey] = upmd
	// 			nd := upmd.UpTime
	// 			if nd <= 0 {
	// 				this.Save()
	// 				atomic.AddInt64(&mod.savenum, 1)
	// 				tk.Stop()
	// 			} else {
	// 				ndt := time.Now().Add(nd)
	// 				if dt.Unix() < time.Now().Unix() {
	// 					//比当前时间小，说明可能是停下来了要重新设置时间
	// 					dt = ndt
	// 					tk.Reset(nd)
	// 				} else if dt.Unix() > ndt.Unix() {
	// 					dt = ndt
	// 					tk.Reset(nd)
	// 				}
	// 			}
	// 		}
	// 	}, func(err interface{}) {
	// 		//如果运行出错了，就等10秒再试
	// 		loglogic.PFatal(err)
	// 		nd := time.Second * 10
	// 		tk.Reset(nd)
	// 		dt = time.Now().Add(nd)
	// 	}, nil)

	// }

}

//Save 执行保存
func (this *DataThread) Save() {
	//保存数据
	for _, data := range this.updatamap {
		if err := data.UpDataSave(this.Conndb); err != nil {
			loglogic.PError(err, " Data: %v ", data)
		}
	}
	this.updatamap = make(map[string]event.ISqlDataModel)
}

// //DataBaseModule 数据持久化模块
// type DataBaseModule struct {
// 	PlayerList map[int]*DataThread //启动的用户协程
// 	// maplock    sync.Mutex            //上面那个集合的锁
// 	ctx        context.Context    //启动控制
// 	cancelfunc context.CancelFunc //取消方法
// 	wg         sync.WaitGroup     //等待所有关闭
// 	currnum    int64              //开协程数
// 	savenum    int64              //保存次数
// 	conndb     *sql.DB            //数据库对象

// 	UpDataChan chan *UpDataModel //写入DB的数据流
// }

// func NewDataBaseModule(sqldb *sql.DB) *DataBaseModule {
// 	result := new(DataBaseModule)
// 	result.conndb = sqldb
// 	result.UpDataChan = make(chan *UpDataModel, 100)
// 	return result
// }

// //Init 初始化
// func (mod *DataBaseModule) Init() {
// 	mod.PlayerList = make(map[int]*DataThread)
// 	mod.ctx, mod.cancelfunc = context.WithCancel(context.Background())

// }

// //Start 启动
// func (mod *DataBaseModule) Start() {
// 	threads.GoTry(mod.DBHandle, nil, nil)
// 	loglogic.PStatus("DataBase Start.")
// }

// //Stop 停止
// func (mod *DataBaseModule) Stop() {
// 	close(mod.UpDataChan)

// 	mod.wg.Wait()
// 	loglogic.PStatus("DataBase Stop.")

// }

// //PrintStatus 打印状态
// func (mod *DataBaseModule) PrintStatus() string {
// 	return fmt.Sprintf(
// 		"\r\n\t\tData Module:	%d/%d	(RunNum/SaveNum)",
// 		atomic.AddInt64(&mod.currnum, 0),
// 		atomic.AddInt64(&mod.savenum, 0))
// }

// func (mod *DataBaseModule) getUserThread(keyid int, conndb *sql.DB) *DataThread {
// 	// mod.maplock.Lock()
// 	// defer mod.maplock.Unlock()
// 	result, ok := mod.PlayerList[keyid]
// 	if !ok {
// 		result = NewDataThread(keyid, conndb)
// 		result.ctx, result.cancelfunc = context.WithCancel(mod.ctx)
// 		result.Go(mod)
// 		mod.PlayerList[keyid] = result
// 	}
// 	return result
// }

// func (mod *DataBaseModule) DelUserThread(keyid int) {
// 	upmd := new(UpDataModel)
// 	upmd.KeyID = keyid
// 	upmd.DataKey = DataBase_DEL_USER_THREAD
// 	mod.AddUpDataModel(upmd)
// 	// mod.maplock.Lock()
// 	// defer mod.maplock.Unlock()

// }

// //AddUpDataModel 添加数据的方法
// func (mod *DataBaseModule) AddUpDataModel(upmd *UpDataModel) {
// 	mod.UpDataChan <- upmd
// }

// //DBHandle 收数据的协程
// func (mod *DataBaseModule) DBHandle() {
// 	mod.wg.Add(1)
// 	defer mod.wg.Done()

// 	for upmd := range mod.UpDataChan {
// 		if upmd.DataKey == DataBase_DEL_USER_THREAD {
// 			result, ok := mod.PlayerList[upmd.KeyID]
// 			if ok {
// 				delete(mod.PlayerList, upmd.KeyID)
// 				result.cancelfunc()
// 				close(result.SendChan)
// 			}
// 		} else {
// 			userth := mod.getUserThread(upmd.KeyID, mod.conndb)
// 			userth.SendChan <- upmd
// 		}

// 	}
// 	mod.cancelfunc()
// 	for _, userth := range mod.PlayerList {
// 		close(userth.SendChan)
// 	}
// 	loglogic.PDebug("DataBase closed player thread.")

// }
