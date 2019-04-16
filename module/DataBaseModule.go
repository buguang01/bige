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

type DataDBModel interface {
}

type UpDataSave func(conndb *sql.DB, datamd DataDBModel) error

//UpDataModel 要修改的数据
type UpDataModel struct {
	KeyID       int32         //用户主键
	DataKey     string        //数据表
	UpTime      time.Duration //保存时间
	SaveFun     UpDataSave    //保存方法
	DataDBModel DataDBModel   //要保存的东西
}

//DataThread 用户的数据库协程
type DataThread struct {
	KeyID     int32                   //用户主键
	updatamap map[string]*UpDataModel //缓存要更新的数据
	SendChan  chan *UpDataModel       //收要更新的数据
	Conndb    *sql.DB                 //数据库连接对象

	ctx        context.Context    //启动控制
	cancelfunc context.CancelFunc //取消方法
}

func NewDataThread(keyid int32, conndb *sql.DB) *DataThread {
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
			d := time.Second * 9999999
			tk := time.NewTicker(d)
			tk.Stop()
			dt := time.Now().Add(d)
			for {
				select {
				case <-this.ctx.Done():
					this.Save()
					atomic.AddInt64(&mod.savenum, 1)

					return
				case <-tk.C:
					this.Save()
					atomic.AddInt64(&mod.savenum, 1)
				case d := <-this.SendChan:
					this.updatamap[d.DataKey] = d
					nd := d.UpTime
					if nd == 0 {
						this.Save()
						atomic.AddInt64(&mod.savenum, 1)
					} else {
						ndt := time.Now().Add(nd)
						if dt.Unix() > ndt.Unix() {
							tk.Stop()
							tk = time.NewTicker(nd)
							dt = ndt
						}
					}

				}
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
		if err := data.SaveFun(this.Conndb, data.DataDBModel); err != nil {
			loglogic.PError(err)
		}
	}
	this.updatamap = make(map[string]*UpDataModel)
}

//DataBaseModule 数据持久化模块
type DataBaseModule struct {
	PlayerList map[int32]*DataThread //启动的用户协程
	maplock    sync.Mutex            //上面那个集合的锁
	ctx        context.Context       //启动控制
	cancelfunc context.CancelFunc    //取消方法
	wg         sync.WaitGroup        //等待所有关闭
	currnum    int64                 //开协程数
	savenum    int64                 //保存次数
}

func NewDataBaseModule() *DataBaseModule {
	result := new(DataBaseModule)
	return result
}

//Init 初始化
func (mod *DataBaseModule) Init() {
	mod.PlayerList = make(map[int32]*DataThread)
	mod.ctx, mod.cancelfunc = context.WithCancel(context.Background())

}

//Start 启动
func (mod *DataBaseModule) Start() {
	loglogic.PStatus("DataBase Start.")
}

//Stop 停止
func (mod *DataBaseModule) Stop() {
	mod.cancelfunc()
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

func (mod *DataBaseModule) GetUserThread(keyid int32, conndb *sql.DB) *DataThread {
	mod.maplock.Lock()
	defer mod.maplock.Unlock()
	result, ok := mod.PlayerList[keyid]
	if !ok {
		result = NewDataThread(keyid, conndb)
		result.ctx, result.cancelfunc = context.WithCancel(mod.ctx)
		result.Go(mod)
		mod.PlayerList[keyid] = result
	}
	return result
}
