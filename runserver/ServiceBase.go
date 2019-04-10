package runserver

import "buguang01/gsframe/module"

//GameConf 游戏服务器的配置
type GameConf struct{
	ServiceID int32//游戏服务器ID
	
}

//GameServiceBase 游戏模块管理
type GameServiceBase struct {
	mlist []*module.IModule
}

func NewGameService() *GameServiceBase {
	result := new(GameServiceBase)
	result.mlist = make([]module, 0, 10) //一般一个服务器能开10个的话就很复杂了
	return result
}

func (gs *GameServiceBase) AddModule(md *module.IModule) {
	gs.mlist = append(gs.mlist, md)
	md.Init()
}

func (gs *GameServiceBase) Run() {
	//
	for _, md := range gs.mlist {
		md.Start()
	}
	//这里要柱塞等关闭
	for _, md := range gs.mlist {
		md.Stop()
	}
}
