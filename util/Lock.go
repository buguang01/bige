package util

import (
	"sync"
)

//UsingRead 读锁
func UsingRead(lk *sync.RWMutex, f func()) {
	lk.RLock()
	defer lk.RUnlock()
	f()
}

//UsingWiter 写锁
func UsingWiter(lk *sync.RWMutex, f func()) {
	lk.Lock()
	defer lk.Unlock()
	f()
}
