package event

import (
	"buguang01/gsframe/loglogic"
	"encoding/json"
	"fmt"
)

//MsgEventBase 消息JSON的基础信息
type MsgEventBase struct {
	Action    int32  //消息号
	ActionKey int32  //消息序号
	MemberID  int32  //用户ID
	Hash      string //发回给用户信息用的钥匙
}

//IMsgEvent 所有的MsgEventBase都要实现这个
type IMsgEvent interface {
	GetMsgHandle()       //刚拿到消息的协程上的调用
	CallbackHandle()     //一般是放到业务成的线程时的回调
	GetAction() int32    //消息号
	GetActionKey() int32 //消息序号
	GetMemberID() int32  //用户ID
	GetHash() string     //发回给用户信息用的钥匙
}

//GetMsgHandle 对IMsgEvent的实现
func (et *MsgEventBase) GetMsgHandle() {
	loglogic.PError(fmt.Sprintf("%T 没有实现GetMsgHandle", et))
}

//CallbackHandle 对IMsgEvent的实现
func (et *MsgEventBase) CallbackHandle() {
	loglogic.PError(fmt.Sprintf("%T 没有实现CallbackHandle", et))
}

//GetAction 消息号
func (et *MsgEventBase) GetAction() int32 {
	return et.Action
}

//GetActionKey 消息序号
func (et *MsgEventBase) GetActionKey() int32 {
	return et.ActionKey
}

//GetMemberID 用户ID
func (et *MsgEventBase) GetMemberID() int32 {
	return et.MemberID

}

//GetHash 发回给用户信息用的钥匙
func (et *MsgEventBase) GetHash() string {

	return et.Hash
}

//IHTTPMsgEVent 如果是HTTP的收到的事件需要这个这个接口
//因为http的请求需要直接回复
type IHTTPMsgEVent interface {
	IMsgEvent
	HTTPGetMsgHandle() <-chan []byte //在HTTP的协程上调用的方法，返回一个在处理完后返回到这个协程的信道
}

//HTTPGetMsgHandle HTTP的协程刚拿到消息的调用
func (et *MsgEventBase) HTTPGetMsgHandle() <-chan []byte {
	result := make(chan []byte, 1)

	//这是一段例子，也是我自己定义的标准回复信息
	resultjs := make(map[string]interface{})
	resultjs["ACTION"] = et.Action
	resultjs["ACTIONCOM"] = 0
	resultjs["ACTIONKEY"] = et.ActionKey
	resultb, _ := json.Marshal(resultjs)
	result <- resultb
	return result
}
