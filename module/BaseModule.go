package module

//IModule 模块接口
type IModule interface {
	//Init 初始化
	Init()
	//Start 启动
	Start()
	//Stop 停止
	Stop()
}
