package runserver

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/module"
)

//GameConfigModel 游戏服务器的配置
type GameConfigModel struct {
	ServiceID   int           //游戏服务器ID
	PStatusTime time.Duration //打印状态的时间（秒）
}

//GameServiceBase 游戏模块管理
type GameServiceBase struct {
	mlist             []module.IModule
	cg                *GameConfigModel
	isrun             bool
	ServiceStopHander func() //当服务器被关掉的时候，先调用的方法
}

//NewGameService 生成一个新的游戏服务器
func NewGameService(conf *GameConfigModel) *GameServiceBase {
	result := new(GameServiceBase)
	result.mlist = make([]module.IModule, 0, 10) //一般一个服务器能开10个的话就很复杂了
	result.cg = conf
	return result
}

//AddModule 给这个管理器，加新的模块
func (gs *GameServiceBase) AddModule(md module.IModule) {
	gs.mlist = append(gs.mlist, md)
	md.Init()
}

//Run 运行游戏
func (gs *GameServiceBase) Run() {
	gs.isrun = true
	//
	for _, md := range gs.mlist {
		md.Start()
	}
	//这里要柱塞等关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

Pstatus:
	for {
		t := time.NewTicker(gs.cg.PStatusTime * time.Second)
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
func (gs *GameServiceBase) GetIsRun() bool {
	return gs.isrun
}
