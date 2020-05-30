package messages

import (
	"bytes"
	"encoding/binary"
)

//路由消息接口
type IGateMessage interface {
	// GetMsgID() uint32
	// GetMyID() uint32
	// GetTargetID() uint32

	// SetMsgID(msgid uint32)
	// SetMyID(myid uint32)
	// SetTargetID(targetid uint32)
	//编码，传出编码的数据和数据的长度
	GateMarshal() ([]byte, uint32)
	//解码,传入数据，传出使用后剩下的数据，和使用了多少字节
	GateUnmarshal(buff []byte) ([]byte, uint32)
}

/*

 */
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

//编码
func (msg *GateMessage) GateMarshal() ([]byte, uint32) {
	buff := &bytes.Buffer{}
	tmpbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(tmpbuf, msg.MsgID)
	buff.Write(tmpbuf)
	binary.BigEndian.PutUint32(tmpbuf, msg.MyID)
	buff.Write(tmpbuf)
	binary.BigEndian.PutUint32(tmpbuf, msg.TargetID)
	buff.Write(tmpbuf)
	return buff.Bytes(), uint32(buff.Len())
}

//解码
func (msg *GateMessage) GateUnmarshal(buff []byte) ([]byte, uint32) {
	msg.SetMsgID(binary.BigEndian.Uint32(buff[:4]))
	buff = buff[4:]
	msg.SetMyID(binary.BigEndian.Uint32(buff[:4]))
	buff = buff[4:]
	msg.SetTargetID(binary.BigEndian.Uint32(buff[:4]))
	buff = buff[4:]
	return buff, 12
}

type DefGateMessage struct {
	GateMessage
}

/*
gate的逻辑是
识别客户端发的消息号，来决定消息要发给哪个服务器，
所以，对客户端业说，没有网关；网关是直转消息；
但是服务器可能不知道这是哪个用户发来的消息；
避免客户端伪造别人发消息，所以消息中用户信息部分应该是gate进行填充；
所以，如果是走gate，那服务器收消息部分应该按gate消息解；
解的方式就是把原消息的消息号拿出来，然后按gate的消息重组；使用的是原消息号

因为前置gate不做服务器之间的消息转发，
所以对与服务来说，知道哪个gate过来的消息是客户过来的消息，还是别的服务器过来的消息



服务器如果过网关回复给客户端的直回消息，使用网关消息ID：XXX（具体看gate消息定义）直回消息
在目标那一项中写上用户的ID就可以了；当然为了减少gate的拆解包的操作，gate应该是在收到后，
直接把后面的数据直发给客户端，所以对与客户端要解的数据的封包操作还是game做的;
对gate是黑盒的；
如果使用服务器消息转发逻辑
一、1对1的服务器转发
服务器发的消息要找网关逻辑来决定这个消息从哪来的，发到哪里去
*/
