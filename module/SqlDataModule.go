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
	playerlist map[int]*DataThread        //操作数据库的列表
	keylist    []int                      //key列表，用来间隔遍历
	chandata   chan []event.ISqlDataModel //消息信道
	mgGo       *threads.ThreadGo          //子协程管理器
	getnum     int64                      //收到的消息数
	conndb     *sql.DB                    //数据库对象
	cg         *SqlDataConfig
}

func NewSqlDataModule(config *SqlDataConfig, sqldb *sql.DB) *SqlDataModule {
	result := new(SqlDataModule)
	result.cg = config
	result.playerlist = make(map[int]*DataThread, config.InitNum)
	result.keylist = make([]int, 0, config.InitNum)
	result.chandata = make(chan []event.ISqlDataModel, config.ChanNum)
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
		"\r\n		Sql Module          :%d/%d		(RunNum/getnum)",
		len(this.playerlist),
		atomic.AddInt64(&this.getnum, 0))
}
func (this *SqlDataModule) AddMsg(msg ...event.ISqlDataModel) {
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
					case upmdarr := <-this.chandata:
						{
							upmd := upmdarr[0]
							logicth, ok := this.playerlist[upmd.GetKeyID()]
							if !ok {
								//新开一个协程
								logicth = newDataThread(upmd.GetKeyID(), this.conndb, this.cg.SubChanNum)
								this.playerlist[logicth.KeyID] = logicth
								logicth.Start(this)
							}
							logicth.UpTime = util.GetCurrTime()
							logicth.SendChan <- upmdarr
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
		case upmdarr, ok := <-this.chandata:
			{
				if !ok {
					return
				}
				atomic.AddInt64(&this.getnum, 1)
				upmd := upmdarr[0]
				logicth, ok := this.playerlist[upmd.GetKeyID()]
				if !ok {
					//新开一个协程
					logicth = newDataThread(upmd.GetKeyID(), this.conndb, this.cg.SubChanNum)
					this.playerlist[logicth.KeyID] = logicth
					this.keylist = append(this.keylist, logicth.KeyID)
					logicth.Start(this)
				}
				logicth.UpTime = util.GetCurrTime()
				logicth.SendChan <- upmdarr
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
	SendChan  chan []event.ISqlDataModel     //收要更新的数据
	Conndb    *sql.DB                        //数据库连接对象
	UpTime    time.Time                      //更新时间
}

func newDataThread(keyid int, conndb *sql.DB, channum int) *DataThread {
	result := new(DataThread)
	result.KeyID = keyid
	result.updatamap = make(map[string]event.ISqlDataModel)
	result.SendChan = make(chan []event.ISqlDataModel, channum)
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
				for upmdarr := range this.SendChan {
					for _, upmd := range upmdarr {
						this.updatamap[upmd.GetDataKey()] = upmd
					}
				}
				threads.Try(this.Save, nil, nil)
				break threadhandle
			}
		case upmdarr, ok := <-this.SendChan:
			{
				if !ok {
					threads.Try(this.Save, nil, nil)
					break threadhandle
				}
				nd := time.Hour
				for _, upmd := range upmdarr {
					this.updatamap[upmd.GetDataKey()] = upmd
					if nd > upmd.GetUpTime() {
						nd = upmd.GetUpTime()
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
				if nd <= 0 {
					threads.Try(this.Save, nil, nil)
					// atomic.AddInt64(&mod.savenum, 1)
					tk.Stop()
				}

			}
		case <-tk.C:
			{
				threads.Try(this.Save, nil, nil)
				tk.Stop()
			}
		}
	}
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
