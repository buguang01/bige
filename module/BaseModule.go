package module

//IModule 模块接口
type IModule interface {
	//Init 初始化
	//configpath:JSON配置文件地址
	Init(configpath string)
	//Start 启动
	Start()
	//Stop 停止
	Stop()
}
