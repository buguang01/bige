// package event

// //INsqdMessage 消息接口
// type INsqdMessage interface {
// 	GetSendID() int
// 	GetSendSID() string
// 	SetSendSID(sid string)
// 	GetActionID() int
// 	GetData() interface{}
// 	GetTopic() string
// }

// type NsqdHander func(msg INsqdMessage)

// type NsqdMessage struct {
// 	SendID   int         //发信息用户ID
// 	SendSID  string      //发信息服务器（回复用的信息）
// 	ActionID int         //消息号
// 	Data     interface{} //消息数据
// 	Topic    string      //目标
// }

// func NewNsqdMessage(mid, actid int, topic string, data interface{}) INsqdMessage {
// 	result := new(NsqdMessage)
// 	result.SendID = mid
// 	result.ActionID = actid
// 	result.Topic = topic
// 	result.Data = data
// 	return result
// }

// func (this *NsqdMessage) GetSendID() int {
// 	return this.SendID
// }
// func (this *NsqdMessage) GetSendSID() string {
// 	return this.SendSID
// }
// func (this *NsqdMessage) SetSendSID(sid string) {
// 	this.SendSID = sid
// }
// func (this *NsqdMessage) GetActionID() int {
// 	return this.ActionID
// }
// func (this *NsqdMessage) GetData() interface{} {
// 	return this.Data
// }
// func (this *NsqdMessage) GetTopic() string {
// 	return this.Topic
// }
