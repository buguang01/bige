package messages

//路由消息接口
type IGateMessage interface {
	GetMsgID() uint32
	GetMyID() uint32
	GetTargetID() uint32

	SetMsgID(msgid uint32)
	SetMyID(myid uint32)
	SetTargetID(targetid uint32)
}

type GateMessage struct {
	MsgID    uint32 `json:"-"` //消息号
	MyID     uint32 `json:"-"` //消息源ID
	TargetID uint32 `json:"-"` //目标ID
}

func (msg *GateMessage) GetMsgID() uint32 {
	return msg.MsgID
}
func (msg *GateMessage) GetMyID() uint32 {
	return msg.MyID
}
func (msg *GateMessage) GetTargetID() uint32 {
	return msg.TargetID
}
func (msg *GateMessage) SetMsgID(msgid uint32) {
	msg.MsgID = msgid
}
func (msg *GateMessage) SetMyID(myid uint32) {
	msg.MyID = myid
}
func (msg *GateMessage) SetTargetID(targetid uint32) {
	msg.TargetID = targetid
}

type DefGateMessage struct {
	GateMessage
}
