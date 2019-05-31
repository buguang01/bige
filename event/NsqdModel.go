package event

type NsqdHander func(msg *NsqdMessage)

type NsqdMessage struct {
	SendID   int         //发信息用户ID
	SendSID  string      //发信息服务器（回复用的信息）
	ActionID int         //消息号
	Data     interface{} //消息数据
	Topic    string      //目标
}

func NewNsqdMessage(mid, actid int, topic string, data interface{}) *NsqdMessage {
	result := new(NsqdMessage)
	result.SendID = mid
	result.ActionID = actid
	result.Topic = topic
	result.Data = data
	return result
}
