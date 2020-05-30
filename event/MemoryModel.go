// package event

// //IMemoryModel
// //
// //
// //RunAutoEvents\UnloadRun\DoneRun
// //这三个方法是在同一个协程上运行的
// //
// //
// type IMemoryModel interface {
// 	//放入管理器的KEY
// 	GetKey() string
// 	//确认加入到了管理器中后，用来开启，这个数据的一些自动任务
// 	//如果用这个方法本自来启动任务，就可以用对应的这些方法来关闭自动任务
// 	RunAutoEvents()
// 	//时间到时，运行的方法,如果发出了委托，就返回true
// 	UnloadRun() bool
// 	//当服务关闭时，运行的方法，这个时候可能就不清内存了，只是关一些自动任务
// 	DoneRun()
// }

// type MemoryMsg struct {
// 	IMemoryModel
// 	CmdType int //消息命令类型，1：添加；2：删除；
// }
