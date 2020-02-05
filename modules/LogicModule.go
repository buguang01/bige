package modules

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util"

	"github.com/buguang01/util/threads"
)

type LogicModule struct {
	channum   int                         //通道缓存空间
	timeout   time.Duration               //超时时间
	logicList map[int]*logicThread        //子逻辑列表
	keylist   []int                       //key列表，用来间隔遍历
	chanLogic chan messages.ILogicMessage //消息信道
	getnum    int64                       //收到的总消息数
	endnum    int64                       //处理结束数
	runing    int64                       //当前处理的消息
	thgo      *threads.ThreadGo           //子协程管理器
}

func NewLogicModule(opts ...options) *LogicModule {
	result := &LogicModule{
		channum:   1024,
		timeout:   60 * time.Second,
		logicList: make(map[int]*logicThread, runtime.NumCPU()*10),
		keylist:   make([]int, 0, runtime.NumCPU()*10),
		getnum:    0,
		endnum:    0,
		runing:    0,
		thgo:      threads.NewThreadGo(),
	}
	return result
}

func (mod *LogicModule) Init() {
	mod.chanLogic = make(chan messages.ILogicMessage, mod.channum)
}

func (mod *LogicModule) Start() {
	mod.thgo.Go(mod.Hander)
	Logger.PStatus("Logic Module Start!")
}

func (mod *LogicModule) Stop() {
	//但是子协程很有可能会再发消息出来走逻辑。这是要注意的点
	// close(this.chanLogic) //可以在这里关是因为到这一步的时候，可以认为不会再有其他模块发东西过来了。
	mod.thgo.CloseWait()
	Logger.PStatus("Logic Module Shop!")
}

//PrintStatus IModule 接口实现，打印状态
func (mod *LogicModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tLogic Module\t:%d/%d/%d/%d\t(logicnum/get/end/run)",
		len(mod.logicList),
		atomic.LoadInt64(&mod.getnum),
		atomic.LoadInt64(&mod.endnum),
		atomic.LoadInt64(&mod.runing))
}

//AddMsg 发消息给逻辑协程处理
func (mod *LogicModule) AddMsg(logicmd messages.ILogicMessage) {

	atomic.AddInt64(&mod.getnum, 1)
	mod.chanLogic <- logicmd

}

func (mod *LogicModule) Hander(ctx context.Context) {
	tk := time.NewTicker(1 * time.Second)
	loop := 0
	for {
		select {
		case <-ctx.Done():
			{
				/*
					判断是不是关闭服务了，
					如果关闭服务了，那就等处理是不是完成，
					如果处理也都完成了，那就结束协程
				*/
				if mod.getnum == mod.endnum {
					return
				}
			}
		case msg := <-mod.chanLogic:
			{
				lth, ok := mod.logicList[msg.TheardID()]
				if !ok {
					//新开一个协程
					lth = newLogicThread(msg.TheardID(), mod.channum)
					mod.logicList[lth.ThreadID] = lth
					mod.keylist = append(mod.keylist, lth.ThreadID)
					lth.start(mod)
				}
				//收到消息，发给子模块
				lth.addMsg(msg)
			}
		case <-tk.C:
			{
				lilen := len(mod.keylist)
				if lilen == 0 {
					break
				}
				loop = loop % lilen
				keyid := mod.keylist[loop]
				if lth, ok := mod.logicList[keyid]; ok {
					if lth.GetMsgNum() == 0 &&
						util.GetCurrTime().Sub(lth.UpTime) > mod.timeout {
						//确定子协程可以关闭
						lth.stop()
					}
				}
				loop++
			}
		}
	}
}

//LogicThread 逻辑协程
type logicThread struct {
	ThreadID  int                         //协程key
	chanLogic chan messages.ILogicMessage //要处理的逻辑
	UpTime    time.Time                   //更新时间
	cancel    context.CancelFunc          //关闭
}

func newLogicThread(id, channum int) *logicThread {
	result := &logicThread{
		ThreadID:  id,
		chanLogic: make(chan messages.ILogicMessage, channum),
		UpTime:    util.GetCurrTime(),
	}
	return result
}

func (lth *logicThread) start(mod *LogicModule) {
	lth.cancel = mod.thgo.SubGo(
		func(ctx context.Context) {
			lth.handle(ctx, mod)
		},
	)
}

func (lth *logicThread) stop() {
	lth.cancel()
	close(lth.chanLogic)
}

func (lth *logicThread) handle(ctx context.Context, mod *LogicModule) {

trheadhandle:
	for {
		select {
		case msg, ok := <-lth.chanLogic:
			{
				if !ok {
					break trheadhandle
				}
				atomic.AddInt64(&mod.runing, 1)
				threads.Try(msg.MessageHandle, nil, nil)
				atomic.AddInt64(&mod.runing, -1)
				atomic.AddInt64(&mod.endnum, 1)
			}
		}
	}
}

func (lth *logicThread) addMsg(msg messages.ILogicMessage) {
	lth.UpTime = util.GetCurrTime()
	lth.chanLogic <- msg
}

//还有多少消息没有处理完
func (lth *logicThread) GetMsgNum() int {
	return len(lth.chanLogic)
}
