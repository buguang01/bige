package runserver

import (
	"buguang01/gsframe/module"
	"os"
	"os/signal"
	"time"
)

//GameConfigModel 游戏服务器的配置
type GameConfigModel struct {
	ServiceID   int32         //游戏服务器ID
	PStatusTime time.Duration //打印状态的时间（秒）
}

//GameServiceBase 游戏模块管理
type GameServiceBase struct {
	mlist []module.IModule
	cg    *GameConfigModel
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
	//
	for _, md := range gs.mlist {
		md.Start()
	}
	//这里要柱塞等关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

Pstatus:
	for {
		t := time.NewTicker(gs.cg.PStatusTime * time.Second)
		select {
		case <-c: //退出
			break Pstatus
		case <-t.C:
			for _, md := range gs.mlist {
				md.PrintStatus()
			}
		}
	}
	for i := len(gs.mlist) - 1; i >= 0; i-- {
		md := gs.mlist[i]
		md.Stop()
	}

}
