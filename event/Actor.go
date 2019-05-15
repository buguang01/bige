package event

type ActorBase struct {
	EventStepNum int //当前步骤
}

//IActor 消息盒子的借口
type IActor interface {
	//当前要运行的步骤
	GetStepNum() int
	//运行步骤成功就进入下一步，失败就开始执行回滚，成功失败内部处理
	RunNext()
	//是否已完成
	IsEnd() bool
	//拿到要返回用的JSON对象，或不用这个直接推信息
	GetResult() JsonMap


}
