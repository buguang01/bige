package threads

import (
	"buguang01/gsframe/loglogic"
	"context"
	"sync"
	"time"
)

//GoTry 在新线程上跑
func GoTry(f func(), catch func(interface{}), finally func()) {
	go Try(f, catch, finally)
}

//Try C#中 的try
func Try(f func(), catch func(interface{}), finally func()) {
	defer func() {
		if finally != nil {
			finally()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			if catch != nil {
				catch(err)
			} else {
				loglogic.PFatal(err)
			}
		}
	}()
	f()

}

//ThreadRun 新开协程的类有回调用的
type ThreadRun struct {
	Chanresult chan func()
}

//NewGoRun 开一个新的协程并运行它
//在新协程上调用f ，resultf回到主线程的方法
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
	GoTry(f, nil, func() {
		this.Chanresult <- resultf
		close(this.Chanresult)
	})
}

//ThreadGo 子协程管理计数，可以等子协程都完成
//用它来管理所有开的协程，需要等这些线程都跑完
type ThreadGo struct {
	Wg  sync.WaitGroup //等待
	Ctx context.Context
	Cal context.CancelFunc
}

func NewThreadGo() *ThreadGo {
	reuslt := new(ThreadGo)
	reuslt.Ctx, reuslt.Cal = context.WithCancel(context.Background())
	return reuslt

}

func (this *ThreadGo) CloseWait() {
	this.Cal()
	this.Wg.Wait()
}

//Go 在当前线程上跑
func (this *ThreadGo) Go(f func(ctx context.Context)) {
	this.Wg.Add(1)
	GoTry(func() {
		f(this.Ctx)
	}, nil, func() {
		this.Wg.Done()
	})
}

//GoTry 在新协程上跑
func (this *ThreadGo) GoTry(f func(ctx context.Context), catch func(interface{}), finally func()) {
	this.Wg.Add(1)
	GoTry(
		func() {
			f(this.Ctx)
		},
		catch,
		func() {
			defer this.Wg.Done()
			if finally != nil {
				finally()
			}
		})
}

//Try 在当前协程上跑
func (this *ThreadGo) Try(f func(ctx context.Context), catch func(interface{}), finally func()) {
	this.Wg.Add(1)
	Try(
		func() {
			f(this.Ctx)
		},
		catch,
		func() {
			defer this.Wg.Done()
			if finally != nil {
				finally()
			}
		})
}

func TimeoutGo(f func(), ticker time.Ticker, timeoutfunc func()) {
	result := make(chan struct{})

	GoTry(f, nil, func() {
		result <- struct{}{}
		close(result)
	})
	select {
	case <-result:

	case <-ticker.C:
		//写一些一定要打印的信息
		if timeoutfunc != nil {
			timeoutfunc()
		}
	}
}
