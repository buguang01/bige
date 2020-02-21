package modules

import (
	"context"
	"fmt"
	"sync"

	"github.com/buguang01/Logger"
	"github.com/buguang01/util/threads"
)

/*
启动任务会开子协程进行，关闭的时候使用,会等协程停下来才会继续
任务停止会占用锁，建议不要有大逻辑，而是发消息出去处理
*/

type AutoTaskModule struct {
	taskList map[string]IAutoTaskModel //任务列表
	lock     *sync.Mutex               //锁
	thgo     *threads.ThreadGo         //协程管理器
}

func NewAutoTaskModule(opts ...options) *AutoTaskModule {
	result := &AutoTaskModule{
		taskList: make(map[string]IAutoTaskModel),
		lock:     &sync.Mutex{},
		thgo:     threads.NewThreadGo(),
	}
	return result
}

//Init 初始化
func (mod *AutoTaskModule) Init() {

}

//Start 启动
func (mod *AutoTaskModule) Start() {

	Logger.PStatus("AutoTask Module Start.")

}

//Stop 停止
func (mod *AutoTaskModule) Stop() {
	mod.thgo.CloseWait()
	Logger.PStatus("AutoTask Module Stop.")

}

//PrintStatus 打印状态
func (mod *AutoTaskModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tAuto Module\t:%d\t(tasknum)",
		len(mod.taskList))
}

//添加新任务，如果任务已存在，就返回false
func (mod *AutoTaskModule) AddTask(task IAutoTaskModel) bool {
	mod.lock.Lock()
	defer mod.lock.Unlock()
	if _, ok := mod.taskList[task.GetTaskName()]; ok {
		return false
	} else {
		mod.taskList[task.GetTaskName()] = task
		task.Start(mod)
		return true
	}
}

//添加新任务，如果任务已存在，就把旧的给关掉
func (mod *AutoTaskModule) ReTask(task IAutoTaskModel) bool {
	mod.lock.Lock()
	defer mod.lock.Unlock()
	if tk, ok := mod.taskList[task.GetTaskName()]; ok {
		tk.Stop()
		mod.taskList[task.GetTaskName()] = task
		task.Start(mod)
		return true
	} else {
		mod.taskList[task.GetTaskName()] = task
		task.Start(mod)
		return true
	}
}

//停止指定任务
func (mod *AutoTaskModule) DelRask(name string) bool {
	mod.lock.Lock()
	defer mod.lock.Unlock()
	if tk, ok := mod.taskList[name]; ok {
		tk.Stop()
		delete(mod.taskList, name)
		return true
	}
	return true
}

/*
循环任务接口
使用AddTask添加到module中时，如果任务以存在就不做任何事
使用ReTask添加到module中时，如果任务已存在就会先把之前的停下来，再把自己添加到管理器中
*/
type IAutoTaskModel interface {
	//任务名字（唯一性）
	GetTaskName() string
	//开始任务
	Start(mod *AutoTaskModule)
	//结束任务
	Stop()
}

type AutoTaskModel struct {
	thgo   *threads.ThreadGo
	Handle func(ctx context.Context) //需要实现这个方法
}

//任务名字（唯一性）需要重载更新这个名字
func (task *AutoTaskModel) GetTaskName() string {
	return "auto"
}

//开始任务
func (task *AutoTaskModel) Start(mod *AutoTaskModule) {
	task.thgo = threads.NewThreadGoByGo(mod.thgo)
	task.thgo.Go(task.Handle)
}

//结束任务
func (task *AutoTaskModel) Stop() {
	task.thgo.CloseWait()
}

// func (task *AutoTaskModel) Handle(ctx context.Context) {
// 	panic(errors.New("任务需要重载这个方法"))
// }
