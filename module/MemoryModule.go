package module

/**
这个模块用来管理内存数据，什么时候可以释放的管理器
比如，一个用户离线后，一定时间了后，就可以卸载这个用户的数据

MemoryModule.AddListenMsg
是用来添加数据到管理器，也是用来重置时间的方法
数据继承IMemoryModel接口，实现里面的所有方法
其中GetKey用来表示数据的唯一标识
RunAutoEvents，表示添加进管理器时运行的方法，如果之前已在管理器中还没有被去除就不会运行它
UnloadRun，设置的时间到了之后运行的方法，如果你确认这个数据是要被卸载的话，就返回true
DoneRun,当服务器要关闭的时候，会运行的方法。可以用于关闭在RunAutoEvents里启动的协程
以上三个方法都在同一个协程运行。

*/
import (
	"buguang01/gsframe/event"
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/threads"
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

//MemoryConfig 内存模块配置
type MemoryConfig struct {
	RunTime int //空闲时间（秒）当这个数据这么长时间 没有被访问时就可以运行了
	InitNum int //初始化内存空间
	ChanNum int //通道缓存空间
}

//MemoryModule 内存缓存管理器
type MemoryModule struct {
	mdList    map[string]*MemoryThread //需要管理的内存单元
	chandata  chan *event.MemoryMsg    //放入管理单元
	mgGo      *threads.ThreadGo        //子协程管理器
	loadNum   int64                    //加载数
	unloadNum int64                    //卸载数
	cg        *MemoryConfig            //配置
}

func NewMemoryModule(config *MemoryConfig) *MemoryModule {
	result := new(MemoryModule)
	result.cg = config
	result.mdList = make(map[string]*MemoryThread, config.InitNum)
	result.chandata = make(chan *event.MemoryMsg, config.ChanNum)
	result.mgGo = threads.NewThreadGo()
	return result
}

//Init 初始化
func (this *MemoryModule) Init() {
	this.loadNum = 0
	this.unloadNum = 0
}

//Start 启动
func (this *MemoryModule) Start() {
	this.mgGo.Go(this.Handle)
	loglogic.PStatus("Momery Module Start!")
}

//Stop 停止
func (this *MemoryModule) Stop() {
	this.mgGo.CloseWait()
	loglogic.PStatus("Momery Module Start!")
}

//PrintStatus 打印状态
func (this *MemoryModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n		Momery Module       :%d/%d/%d	(sum/load/unload)",
		len(this.mdList),
		atomic.AddInt64(&this.loadNum, 0),
		atomic.AddInt64(&this.unloadNum, 0))
}

func (this *MemoryModule) AddListenMsg(data event.IMemoryModel) {
	select {
	case <-this.mgGo.Ctx.Done():
	default:
		{
			msg := new(event.MemoryMsg)
			msg.CmdType = 1
			msg.IMemoryModel = data
			this.chandata <- msg
		}
	}

}

//
func (this *MemoryModule) reListenMsg(data event.IMemoryModel) {
	select {
	case <-this.mgGo.Ctx.Done():
	default:
		{
			msg := new(event.MemoryMsg)
			msg.CmdType = 2
			msg.IMemoryModel = data
			this.chandata <- msg
		}
	}
}

func (this *MemoryModule) Handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case cmd := <-this.chandata:
			{
				//收到数据
				key := cmd.GetKey()
				switch cmd.CmdType {
				case 1: //添加
					{
						th, ok := this.mdList[key]
						if ok {
							th.tk.Reset(time.Duration(this.cg.RunTime) * time.Second)
						} else {
							th = new(MemoryThread)
							th.Key = key
							th.tk = time.NewTimer(time.Duration(this.cg.RunTime) * time.Second)
							th.DataMemory = cmd.IMemoryModel
							this.mdList[key] = th
							//启动
							this.mgGo.Go(func(ctx context.Context) {
								th.Handle(ctx, this)
							})
							atomic.AddInt64(&this.loadNum, 1)
						}
					}
				case 2: //删除
					{
						_, ok := this.mdList[key]
						if ok {
							delete(this.mdList, key)
							atomic.AddInt64(&this.unloadNum, 1)
						}
					}
				}
			}
		}
	}
}

type MemoryThread struct {
	Key        string      //主键
	tk         *time.Timer //计时器
	DataMemory event.IMemoryModel
}

func (this *MemoryThread) Handle(ctx context.Context, mg *MemoryModule) {
	threads.Try(this.DataMemory.RunAutoEvents, nil, nil)
memorythread:
	for {
		select {
		case <-ctx.Done():
			{
				//管理器要关闭了
				threads.Try(this.DataMemory.DoneRun, nil, nil)
				break memorythread
			}
		case <-this.tk.C:
			{
				this.tk.Stop()
				//时间到了，发消息出去
				result := false
				threads.Try(func() {
					result = this.DataMemory.UnloadRun()
				}, nil, nil)

				if result {
					mg.reListenMsg(this.DataMemory)
					break memorythread
				}
				this.tk.Reset(time.Duration(mg.cg.RunTime) * time.Second)
				// this.tk.Stop()
			}
		}
	}

}
