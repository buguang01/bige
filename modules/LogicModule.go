package modules

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util"

	"github.com/buguang01/util/threads"
)

//设置通道缓存空间
func LogicSetChanNum(channum int) options {
	return func(mod IModule) {
		mod.(*LogicModule).chanNum = channum
	}
}

//设置超时时间(秒）
func LogicSetTimeout(timeout time.Duration) options {
	return func(mod IModule) {
		mod.(*LogicModule).timeout = timeout * time.Second
	}
}

type LogicModule struct {
	chanNum   int                         //通道缓存空间
	timeout   time.Duration               //超时时间
	logicList map[int]*logicThread        //子逻辑列表
	keyList   []int                       //key列表，用来间隔遍历
	chanLogic chan messages.ILogicMessage //消息信道
	getNum    int64                       //收到的总消息数
	endNum    int64                       //处理结束数
	runing    int64                       //当前处理的消息
	thgo      *threads.ThreadGo           //子协程管理器
}

func NewLogicModule(opts ...options) *LogicModule {
	result := &LogicModule{
		chanNum:   1024,
		timeout:   60 * time.Second,
		logicList: make(map[int]*logicThread, moduleCap),
		keyList:   make([]int, 0, moduleCap),
		getNum:    0,
		endNum:    0,
		runing:    0,
		thgo:      threads.NewThreadGo(),
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func (mod *LogicModule) Init() {
	mod.chanLogic = make(chan messages.ILogicMessage, mod.chanNum)
}

func (mod *LogicModule) Start() {
	mod.thgo.Go(mod.Handle)
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
		atomic.LoadInt64(&mod.getNum),
		atomic.LoadInt64(&mod.endNum),
		atomic.LoadInt64(&mod.runing))
}

//AddMsg 发消息给逻辑协程处理
func (mod *LogicModule) AddMsg(logicmd messages.ILogicMessage) {

	atomic.AddInt64(&mod.getNum, 1)
	mod.chanLogic <- logicmd

}

func (mod *LogicModule) Handle(ctx context.Context) {
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
					当确定getnum==endnum时就是所有消息都处理完了
				*/
				if atomic.LoadInt64(&mod.getNum) == atomic.LoadInt64(&mod.endNum) {
					for _, lth := range mod.logicList {
						lth.stop()
					}
					return
				}
			}
		case msg, ok := <-mod.chanLogic:
			{
				if !ok {
					continue
				}
				lth, ok := mod.logicList[msg.LogicThreadID()]
				if !ok {
					//新开一个协程
					lth = newLogicThread(msg.LogicThreadID(), mod.chanNum)
					mod.logicList[lth.LogicThreadID] = lth
					mod.keyList = append(mod.keyList, lth.LogicThreadID)
					lth.start(mod)
				}
				//收到消息，发给子模块
				lth.addMsg(msg)
			}
		case <-tk.C:
			{
				lilen := len(mod.keyList)
				if lilen == 0 {
					break
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

//LogicThread 逻辑协程
type logicThread struct {
	LogicThreadID int                         //协程key
	chanLogic     chan messages.ILogicMessage //要处理的逻辑
	upTime        time.Time                   //更新时间
	cancel        context.CancelFunc          //关闭
}

func newLogicThread(id, channum int) *logicThread {
	result := &logicThread{
		LogicThreadID: id,
		chanLogic:     make(chan messages.ILogicMessage, channum),
		upTime:        util.GetCurrTime(),
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
				atomic.AddInt64(&mod.endNum, 1)
			}
		}
	}
}

func (lth *logicThread) addMsg(msg messages.ILogicMessage) {
	lth.upTime = util.GetCurrTime()
	lth.chanLogic <- msg
}

//还有多少消息没有处理完
func (lth *logicThread) GetMsgNum() int {
	return len(lth.chanLogic)
}
