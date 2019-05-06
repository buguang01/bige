package module

import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/threads"
	"buguang01/gsframe/util"
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

//逻辑模块
//一般会用一个ID，来确认这个消息在哪个协程上处理
//这个ID可以是用户ID，也可以是消息，也可以是结合起来使用
//这样做的目地是为了同一用户的同一类消息会在一个线程上处理

type LogicConfig struct {
	Timeout    int //超时时间（秒）
	InitNum    int //初始化内存空间
	ChanNum    int //通道缓存空间
	SubChanNum int //子通道缓存空间
}

//LogicModule 逻辑模块
type LogicModule struct {
	cg        *LogicConfig                  //配置信息
	logicList map[string]*LogicThreadModule //子逻辑列表
	keylist   []string                      //key列表，用来间隔遍历
	chanLogic chan event.LogicModel         //消息信道
	getnum    int64                         //收到的总消息数
	currnum   int64                         //当前处理的消息
	mgGo      *threads.ThreadGo             //子协程管理器
}

func NewLogicModule(config *LogicConfig) *LogicModule {
	result := new(LogicModule)
	result.cg = config
	result.logicList = make(map[string]*LogicThreadModule, config.InitNum)
	result.keylist = make([]string, 0, config.InitNum)
	result.chanLogic = make(chan event.LogicModel, config.ChanNum)
	result.mgGo = threads.NewThreadGo()
	return result
}

func (this *LogicModule) Init() {

}

func (this *LogicModule) Start() {
	this.mgGo.Go(this.Hander)
	loglogic.PStatus("Logic Module Start!")
}

func (this *LogicModule) Stop() {
	//但是子协程很有可能会再发消息出来走逻辑。这是要注意的点
	// close(this.chanLogic) //可以在这里关是因为到这一步的时候，可以认为不会再有其他模块发东西过来了。
	this.mgGo.CloseWait()
	loglogic.PStatus("Logic Module Shop!")
}

//PrintStatus IModule 接口实现，打印状态
func (mod *LogicModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n		Logic Module:	%d/%d/%d	(logicnum/get/run)",
		len(mod.logicList),
		atomic.AddInt64(&mod.getnum, 0),
		atomic.AddInt64(&mod.currnum, 0))
}

//AddMsg 发消息给逻辑协程处理
func (this *LogicModule) AddMsg(logicmd event.LogicModel) {
	atomic.AddInt64(&this.currnum, 1)
	this.chanLogic <- logicmd
}

func (this *LogicModule) Hander(ctx context.Context) {
	tk := time.NewTimer(1 * time.Second)
	loop := 0
	for {
		select {
		case <-ctx.Done():
			{
				for {
					select {
					case logicmd, ok := <-this.chanLogic:
						{
							if !ok {
								return
							}
							atomic.AddInt64(&this.getnum, 1)
							logicth, ok := this.logicList[logicmd.KeyID()]
							if !ok {
								//新开一个协程
								logicth = newLogicThread(logicmd.KeyID(), this.cg.SubChanNum)
								this.logicList[logicth.KeyID] = logicth
								this.keylist = append(this.keylist, logicth.KeyID)
								logicth.Start(this)
							}
							//收到消息，发给子模块
							logicth.UpTime = util.GetCurrTime()
							logicth.chanLogic <- logicmd

						}
					default:
						if this.currnum == 0 {
							close(this.chanLogic)
							for _, logicth := range this.logicList {
								logicth.CloseChan()
							}
							return
						}
					}
				}
				// for logicmd := range this.chanLogic {
				// 	logicth, ok := this.logicList[logicmd.KeyID()]
				// 	if !ok {
				// 		//新开一个协程
				// 		logicth = newLogicThread(logicmd.KeyID())
				// 		this.logicList[logicth.KeyID] = logicth
				// 		logicth.Start(this)
				// 	}
				// 	//收到消息，发给子模块
				// 	logicth.chanLogic <- logicmd
				// 	this.getnum++
				// }

			}
		case logicmd, ok := <-this.chanLogic:
			{
				if !ok {
					for _, logicth := range this.logicList {
						logicth.CloseChan()
					}
					return
				}
				atomic.AddInt64(&this.getnum, 1)
				logicth, ok := this.logicList[logicmd.KeyID()]
				if !ok {
					//新开一个协程
					logicth = newLogicThread(logicmd.KeyID(), this.cg.SubChanNum)
					this.logicList[logicth.KeyID] = logicth
					this.keylist = append(this.keylist, logicth.KeyID)
					logicth.Start(this)
				}
				//收到消息，发给子模块
				logicth.UpTime = util.GetCurrTime()
				atomic.AddInt64(&logicth.currnum, 1)
				logicth.chanLogic <- logicmd
			}
		case <-tk.C:
			{
				tk.Reset(1 * time.Second)
				if len(this.keylist) == 0 {
					break
				}
				loop := (loop + 1) % len(this.keylist)
				keyid := this.keylist[loop]
				logicth, ok := this.logicList[keyid]
				if ok {
					if logicth.UpTime.Add(time.Duration(this.cg.Timeout)*time.Second).Unix() < util.GetCurrTime().Unix() &&
						logicth.currnum == 0 {
						logicth.Stop(this, loop)
						loop--
					}
				} else {
					this.keylist = append(this.keylist[:loop], this.keylist[loop+1:]...)
				}
				// //每次就检查10个
				// i := 0
				// for _, logicth := range this.logicList {
				// 	i++
				// 	if logicth.UpTime.Add(time.Duration(this.cg.Timeout)*time.Second).Unix() < util.GetCurrTime().Unix() &&
				// 		logicth.currnum == 0 {
				// 		logicth.Stop(this)
				// 	}
				// 	// else if !logicth.IsRun {
				// 	// 	//如果停下了就从集合里删除了
				// 	// 	delete(this.logicList, logicth.KeyID)
				// 	// }
				// 	if i > 10 {
				// 		break
				// 	}
				// }
				// tk.Reset(10 * time.Second)
			}
		}
	}
}

//LogicThreadModule 逻辑协程
type LogicThreadModule struct {
	KeyID     string                //协程key
	chanLogic chan event.LogicModel //要处理的逻辑
	UpTime    time.Time             //更新时间
	currnum   int64                 //要处理的消息数
}

func newLogicThread(keyid string, channum int) *LogicThreadModule {
	result := new(LogicThreadModule)
	result.KeyID = keyid
	result.chanLogic = make(chan event.LogicModel, channum)
	return result
}

func (this *LogicThreadModule) Start(mg *LogicModule) {
	this.currnum = 0
	mg.mgGo.Go(
		func(ctx context.Context) {
			this.Handle(ctx, mg)
		})

}
func (this *LogicThreadModule) Stop(mg *LogicModule, index int) {
	delete(mg.logicList, this.KeyID)
	mg.keylist = append(mg.keylist[:index], mg.keylist[index+1:]...)
	close(this.chanLogic)
}
func (this *LogicThreadModule) CloseChan() {
	close(this.chanLogic)
}

func (this *LogicThreadModule) Handle(ctx context.Context, mg *LogicModule) {

trheadhandle:
	for {
		select {
		case <-ctx.Done():
			{
				for logicmd := range this.chanLogic {
					threads.Try(logicmd.Run, nil, nil)
					atomic.AddInt64(&this.currnum, -1)
					atomic.AddInt64(&mg.currnum, -1)
				}
				break trheadhandle
			}
		case logicmd, ok := <-this.chanLogic:
			{
				if !ok {
					break trheadhandle
				}
				threads.Try(logicmd.Run, nil, nil)
				atomic.AddInt64(&this.currnum, -1)
				atomic.AddInt64(&mg.currnum, -1)
			}
		}
	}
}
