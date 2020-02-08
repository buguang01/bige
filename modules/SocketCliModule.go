package modules

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util/threads"
)

/*
这是一个Socket的客户端模块
会简单的在连接断开的时候，重新连接
同步发生消息

*/

//这个连接的名字，比如这个连接的目标是什么就叫什么
func SocketCliSetConnName(name string) options {
	return func(mod IModule) {
		mod.(*SocketCliModule).ConnName = name
	}
}

func SocketCliSetPort(ipport string) options {
	return func(mod IModule) {
		mod.(*SocketCliModule).ipPort = ipport
	}
}

//设置路由
func SocketCliSetRoute(route messages.IMessageHandle) options {
	return func(mod IModule) {
		mod.(*SocketCliModule).RouteHandle = route
	}
}

type SocketCliModule struct {
	ConInfo     interface{}             //自定义的连接信息，给上层逻辑使用
	ConnName    string                  //连接名字
	ipPort      string                  //连接服务器的地址
	RouteHandle messages.IMessageHandle //消息路由
	getnum      int64                   //收到的总消息数
	sendnum     int64                   //发出去的消息数
	conn        net.Conn                //连接
	thgo        *threads.ThreadGo       //协程管理器
	isRun       bool                    //是否运行
}

func NewSocketCliModule(opts ...options) *SocketCliModule {
	result := &SocketCliModule{
		ConnName:    "SocketName",
		ipPort:      ":8082",
		RouteHandle: messages.JsonMessageHandleNew(),
		getnum:      0,
		sendnum:     0,
		thgo:        threads.NewThreadGo(),
		isRun:       false,
	}

	for _, opt := range opts {
		opt(result)
	}
	return result
}

//Init 初始化
func (mod *SocketCliModule) Init() {
	var err error
	mod.conn, err = net.Dial("tcp", mod.ipPort)
	if err != nil {
		panic(err)
	}
}

//Start 启动
func (mod *SocketCliModule) Start() {
	mod.isRun = true
	mod.thgo.Go(mod.hander)
	Logger.PStatus("Socket Cli Module Start.")
}

//Stop 停止
func (mod *SocketCliModule) Stop() {
	mod.isRun = false
	mod.conn.Close()
	mod.thgo.CloseWait()
	Logger.PStatus("Socket Cli Module Stop.")
}

//PrintStatus 打印状态
func (mod *SocketCliModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tsocket cli Module\t:%d/%d\t(get/send)",
		atomic.LoadInt64(&mod.getnum),
		atomic.LoadInt64(&mod.sendnum))
}

func (mod *SocketCliModule) hander(ctx context.Context) {
	buf := &bytes.Buffer{}
	for {
		buff, err := ioutil.ReadAll(mod.conn)
		if err != nil || len(buff) == 0 {
			if err != nil {
				Logger.PDebug("Socket Cli Conn Read Error:%+v.", err)
			} else {
				Logger.PDebug("Socket Cli Conn Read Error EOF.")

			}
			if mod.reConn() {
				continue
			} else {
				return
			}
		}
		buf.Write(buff)
		buff = buf.Bytes()
		msglen, ok := mod.RouteHandle.CheckMaxLenVaild(buff)
		if !ok {
			if msglen == 0 {
				//消息长度异常
				return
			}
			continue
		}
		msg, err := mod.RouteHandle.Unmarshal(buff[:msglen])
		if err != nil {
			Logger.PInfo("socket cli RouteHandle Unmarshal Error:%s", err.Error())
			return
		}
		modmsg, ok := msg.(messages.ISocketMessageHandle)
		if !ok {
			Logger.PInfo("Not is socket cli Msg:%+v", msg)
			return
		} else {
			Logger.PInfo("socket cli Get Msg:%+v", msg)
		}
		buf.Reset()
		if uint32(len(buff)) > msglen {
			buf.Write(buff[msglen:])
		}

		atomic.AddInt64(&mod.getnum, 1)
		mod.thgo.Try(func(ctx context.Context) {
			//因为是主动连接，所以不会返回连接，消息里本身应该会带来源信息
			modmsg.SocketDirectCall(nil)
		}, nil, nil)
	}
}

//如果连接断开重连
func (mod *SocketCliModule) reConn() (result bool) {
	result = false
	if mod.isRun {
		for !result {
			mod.thgo.Try(func(ctx context.Context) {
				mod.conn.Close()
				mod.Init()
				result = true
			}, func(err interface{}) {
				time.Sleep(time.Second)
			}, nil)
		}
	}
	return result
}

//同步写入消息
func (mod *SocketCliModule) AddMsgSyn(msg messages.ISocketResultMessage) error {
	if buff, err := mod.RouteHandle.Marshal(msg.GetAction(), msg); err != nil {
		return err
	} else if _, err = mod.conn.Write(buff); err != nil {
		return err
	}
	atomic.AddInt64(&mod.sendnum, 1)
	return nil
}
