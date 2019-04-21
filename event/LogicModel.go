package event

//LogicModel 逻辑委托
type LogicModel interface {
	KeyID() string //所在协程的KEY
	Run()          //调用方法
}
