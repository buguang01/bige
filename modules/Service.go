package modules

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buguang01/Logger"
)

func GameServiceSetSID(sid int) servOptions {
	return func(mod *GameService) {
		mod.ServiceID = sid
	}
}

func GameServiceSetPTime(ptime time.Duration) servOptions {
	return func(mod *GameService) {
		mod.PStatusTime = ptime * time.Second
	}
}

func GameServiceSetStopHander(hander func()) servOptions {
	return func(mod *GameService) {
		mod.ServiceStopHander = hander
	}
}

type servOptions func(mod *GameService)

type GameService struct {
	ServiceID         int           //游戏服务器ID
	PStatusTime       time.Duration //打印状态的时间（秒）
	mlist             []IModule
	isrun             bool
	ServiceStopHander func() //当服务器被关掉的时候，先调用的方法
}

func NewGameService(opts ...servOptions) *GameService {
	result := &GameService{
		ServiceID:         0,
		PStatusTime:       10 * time.Second,
		mlist:             make([]IModule, 0, 10),
		isrun:             false,
		ServiceStopHander: nil,
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

//AddModule 给这个管理器，加新的模块
func (gs *GameService) AddModule(mds ...IModule) {
	gs.mlist = append(gs.mlist, mds...)
	for _, md := range mds {
		md.Init()
	}
}

//Run 运行游戏
func (gs *GameService) Run() {
	gs.isrun = true
	//
	for _, md := range gs.mlist {
		md.Start()
	}
	//这里要柱塞等关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	t := time.NewTicker(gs.PStatusTime)
	defer t.Stop()
Pstatus:
	for {
		select {
		case <-c: //退出
			break Pstatus
		case <-t.C:
			var ps string
			for _, md := range gs.mlist {
				ps += md.PrintStatus()
			}
			Logger.PStatus(ps)
		}
	}
	gs.isrun = false
	if gs.ServiceStopHander != nil {
		gs.ServiceStopHander()
	}
	for i := len(gs.mlist) - 1; i >= 0; i-- {
		md := gs.mlist[i]
		md.Stop()
	}

}

//GetIsRun 我们游戏是不是还在运行着，如果为false表示我们服务器正在关闭中
func (gs *GameService) GetIsRun() bool {
	return gs.isrun
}
