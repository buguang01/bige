package event

import (
	"buguang01/gsframe/util"
)

//IMsgEvent 所有的MsgEventBase都要实现这个
type IMsgEvent interface {
	GetAction() int32    //消息号
	GetActionKey() int32 //消息序号
	GetMemberID() int32  //用户ID
	GetHash() uint32     //发回给用户信息用的钥匙
}

//IHTTPMsgEVent 如果是HTTP的收到的事件需要这个这个接口
//因为http的请求需要直接回复
// type IHTTPMsgEVent interface {
// 	IMsgEvent
// 	HTTPGetMsgHandle() <-chan []byte //在HTTP的协程上调用的方法，返回一个在处理完后返回到这个协程的信道
// }

//HTTPGetMsgHandle HTTP的协程刚拿到消息的调用
// func (et *MsgEventBase) HTTPGetMsgHandle() <-chan []byte {
// 	result := make(chan []byte, 1)

// 	//这是一段例子，也是我自己定义的标准回复信息
// 	resultjs := make(map[string]interface{})
// 	resultjs["ACTION"] = et.Action
// 	resultjs["ACTIONCOM"] = 0
// 	resultjs["ACTIONKEY"] = et.ActionKey
// 	resultb, _ := json.Marshal(resultjs)
// 	result <- resultb
// 	return result
// }

//JsonMap   收到的JSON数据
type JsonMap map[string]interface{}

//GetAction 消息号
func (js JsonMap) GetAction() int32 {
	return util.Convert.ToInt32(js["ACTION"])
}

//GetActionKey int32 //消息序号
func (js JsonMap) GetActionKey() int32 {
	return util.Convert.ToInt32(js["ACTIONKEY"])
}

//GetMemberID int32  //用户ID
func (js JsonMap) GetMemberID() int32 {
	return util.Convert.ToInt32(js["MEMBERID"])
}

//GetHash string     //发回给用户信息用的钥匙
func (js JsonMap) GetHash() uint32 {
	return uint32(js["HASH"].(float64))
}

//JsonArray JSON数组
type JsonArray []interface{}

func (js JsonArray) Add(d interface{}) {
	js = append(js, d)
}
