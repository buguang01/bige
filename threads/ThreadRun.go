package threads

import (
	"buguang01/gsframe/loglogic"
)

//Try C#中 的try
func Try(f func(), catch func(), finally func()) {
	defer func() {
		if finally != nil {
			finally()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			if catch != nil {
				catch()
			} else {
				loglogic.PFatal(err)
			}
		}
	}()
	f()

}

//ThreadRun 新开协程的类
type ThreadRun struct {
	Chanresult chan func()
}

//NewGoRun 开一个新的协程并运行它
func NewGoRun(f func(), resultf func()) *ThreadRun {
	result := new(ThreadRun)
	result.Chanresult = make(chan func(), 1)
	result.Go(f, resultf)
	return result
}

//NewGo 开一个新的协程对象
func NewGo() *ThreadRun {
	result := new(ThreadRun)
	result.Chanresult = make(chan func(), 1)
	return result

}

//Go 在新协程上调用f ，resultf回到主线程的方法
func (this *ThreadRun) Go(f func(), resultf func()) {
	go func() {
		Try(f, nil, func() {
			this.Chanresult <- resultf
			close(this.Chanresult)
		})
	}()

}
