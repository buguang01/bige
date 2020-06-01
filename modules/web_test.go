package modules_test

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/buguang01/bige/modules"
	"github.com/buguang01/util/threads"
)

var (
	WebmodulesEx *modules.WebModule
)

func TestCtx(t *testing.T) {
	thgo := threads.NewThreadGo()
	thgo.CloseWait()
	for {
		select {
		case <-thgo.Ctx.Done():
			fmt.Println("ctx")
		default:
			fmt.Println("def")
		}
		time.Sleep(time.Second)
	}

}

func TestAddInt(t *testing.T) {
	// ch := make(chan int, 8)
	// ch <- 1
	// ch <- 2
	// ch <- 3
	// close(ch)
	// tk := time.NewTimer(time.Second * 10)
	// for {
	// 	select {
	// 	case <-ch:
	// 		fmt.Println("tk.c then")
	// 		time.Sleep(time.Second)
	// 	case <-tk.C:
	// 		fmt.Println("tk.c ")
	// 		time.Sleep(time.Second)

	// 	}
	// }

	var i int64 = 0
	addint(&i)
	fmt.Println(i)
}

func addint(i *int64) {
	atomic.AddInt64(i, 1)
	defer atomic.AddInt64(i, 10)
	atomic.AddInt64(i, 100)
}

func TestBuffer(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("abcdefg")
	fmt.Println(buf.Bytes())
	tp := buf.Next(3)
	fmt.Println(tp)
	fmt.Println(buf.Bytes())
	buf.WriteString("hijklmn")
	fmt.Println(buf.Bytes())
	fmt.Println(tp)
}
