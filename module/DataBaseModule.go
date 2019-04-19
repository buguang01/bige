package module

import (
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/threads"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DataBase_DEL_USER_THREAD string = "DEL USER THREAD"
)

type DataDBModel interface {
}

type UpDataSave func(conndb *sql.DB, datamd DataDBModel) error

//UpDataModel 要修改的数据
type UpDataModel struct {
	KeyID       int         //用户主键
	DataKey     string        //数据表
	UpTime      time.Duration //保存时间
	SaveFun     UpDataSave    //保存方法
	DataDBModel DataDBModel   //要保存的东西
}

//DataThread 用户的数据库协程
type DataThread struct {
	KeyID     int                   //用户主键
	updatamap map[string]*UpDataModel //缓存要更新的数据
	SendChan  chan *UpDataModel       //收要更新的数据
	Conndb    *sql.DB                 //数据库连接对象

	ctx        context.Context    //启动控制
	cancelfunc context.CancelFunc //取消方法
}

func NewDataThread(keyid int, conndb *sql.DB) *DataThread {
	result := new(DataThread)
	result.KeyID = keyid
	result.updatamap = make(map[string]*UpDataModel)
	result.SendChan = make(chan *UpDataModel, 10)
	result.Conndb = conndb

	return result
}

//Go 保存数据的协程
func (this *DataThread) Go(mod *DataBaseModule) {
	mod.wg.Add(1)
	threads.GoTry(
		func() {
			atomic.AddInt64(&mod.currnum, 1)
			d := time.Second * 600
			dt := time.Now().Add(d)
			tk := time.NewTimer(d)
			tk.Stop()
			runing := true
			for runing {
				threads.Try(func() {
					select {
					case <-this.ctx.Done():
						for upmd := range this.SendChan {
							this.updatamap[upmd.DataKey] = upmd
						}
						this.Save()
						atomic.AddInt64(&mod.savenum, 1)
						runing = false
						return
					case <-tk.C:
						this.Save()
						atomic.AddInt64(&mod.savenum, 1)
						tk.Stop()

					case upmd, ok := <-this.SendChan:
						if !ok {
							this.Save()
							atomic.AddInt64(&mod.savenum, 1)
							runing = false
							return
						}
						this.updatamap[upmd.DataKey] = upmd
						nd := upmd.UpTime
						if nd <= 0 {
							this.Save()
							atomic.AddInt64(&mod.savenum, 1)
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
				}, func(err interface{}) {
					//如果运行出错了，就等10秒再试
					loglogic.PFatal(err)
					nd := time.Second * 10
					tk.Reset(nd)
					dt = time.Now().Add(nd)
				}, nil)

			}
		},
		nil,
		func() {
			mod.wg.Done()
			atomic.AddInt64(&mod.currnum, -1)
		},
	)

}

//Save 执行保存
func (this *DataThread) Save() {
	//保存数据
	for _, data := range this.updatamap {
		if data.SaveFun == nil {
			continue
		}
		if err := data.SaveFun(this.Conndb, data.DataDBModel); err != nil {
			loglogic.PError(err)
		}
	}
	this.updatamap = make(map[string]*UpDataModel)
}

//DataBaseModule 数据持久化模块
type DataBaseModule struct {
	PlayerList map[int]*DataThread //启动的用户协程
	// maplock    sync.Mutex            //上面那个集合的锁
	ctx        context.Context    //启动控制
	cancelfunc context.CancelFunc //取消方法
	wg         sync.WaitGroup     //等待所有关闭
	currnum    int64              //开协程数
	savenum    int64              //保存次数
	conndb     *sql.DB            //数据库对象

	UpDataChan chan *UpDataModel //写入DB的数据流
}

func NewDataBaseModule(sqldb *sql.DB) *DataBaseModule {
	result := new(DataBaseModule)
	result.conndb = sqldb
	result.UpDataChan = make(chan *UpDataModel, 100)
	return result
}

//Init 初始化
func (mod *DataBaseModule) Init() {
	mod.PlayerList = make(map[int]*DataThread)
	mod.ctx, mod.cancelfunc = context.WithCancel(context.Background())

}

//Start 启动
func (mod *DataBaseModule) Start() {
	threads.GoTry(mod.DBHandle, nil, nil)
	loglogic.PStatus("DataBase Start.")
}

//Stop 停止
func (mod *DataBaseModule) Stop() {
	close(mod.UpDataChan)

	mod.wg.Wait()
	loglogic.PStatus("DataBase Stop.")

}

//PrintStatus 打印状态
func (mod *DataBaseModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tData Module:	%d/%d	(RunNum/SaveNum)",
		atomic.AddInt64(&mod.currnum, 0),
		atomic.AddInt64(&mod.savenum, 0))
}

func (mod *DataBaseModule) getUserThread(keyid int, conndb *sql.DB) *DataThread {
	// mod.maplock.Lock()
	// defer mod.maplock.Unlock()
	result, ok := mod.PlayerList[keyid]
	if !ok {
		result = NewDataThread(keyid, conndb)
		result.ctx, result.cancelfunc = context.WithCancel(mod.ctx)
		result.Go(mod)
		mod.PlayerList[keyid] = result
	}
	return result
}

func (mod *DataBaseModule) DelUserThread(keyid int) {
	upmd := new(UpDataModel)
	upmd.KeyID = keyid
	upmd.DataKey = DataBase_DEL_USER_THREAD
	mod.AddUpDataModel(upmd)
	// mod.maplock.Lock()
	// defer mod.maplock.Unlock()

}

//AddUpDataModel 添加数据的方法
func (mod *DataBaseModule) AddUpDataModel(upmd *UpDataModel) {
	mod.UpDataChan <- upmd
}

//DBHandle 收数据的协程
func (mod *DataBaseModule) DBHandle() {
	mod.wg.Add(1)
	defer mod.wg.Done()

	for upmd := range mod.UpDataChan {
		if upmd.DataKey == DataBase_DEL_USER_THREAD {
			result, ok := mod.PlayerList[upmd.KeyID]
			if ok {
				delete(mod.PlayerList, upmd.KeyID)
				result.cancelfunc()
				close(result.SendChan)
			}
		} else {
			userth := mod.getUserThread(upmd.KeyID, mod.conndb)
			userth.SendChan <- upmd
		}

	}
	mod.cancelfunc()
	for _, userth := range mod.PlayerList {
		close(userth.SendChan)
	}
	loglogic.PDebug("DataBase closed player thread.")

}
