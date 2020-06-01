package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/buguang01/util"
)

//这个还没有完成
type GateJsonMessageHandle struct {
	msgHead   uint32                 //消息头
	msgMaxLen uint16                 //消息最大长度
	routelist map[uint32]interface{} //消息路由列表
	defType   interface{}            //消息默认类型如果设置了，会在对没有处理的消息进行处理
}

func GateJsonMessageSetMsgHead(msghead uint32) options {
	return func(msghandle IMessageHandle) {
		msghandle.(*GateJsonMessageHandle).msgHead = msghead
	}
}

func GateJsonMessageSetMsgMaxLen(max uint16) options {
	return func(msghandle IMessageHandle) {
		msghandle.(*GateJsonMessageHandle).msgMaxLen = max
	}
}

func GateJsonMessageHandleNew(opts ...options) (msghandle *GateJsonMessageHandle) {
	msghandle = &GateJsonMessageHandle{
		routelist: make(map[uint32]interface{}),
		msgHead:   uint32(0x12340000),
		msgMaxLen: ^uint16(0),
	}
	for _, opt := range opts {
		opt(msghandle)
	}
	return msghandle
}

func (msghandle *GateJsonMessageHandle) GateMarshal(gate IGateMessage, data interface{}) ([]byte, error) {
	buff := &bytes.Buffer{}
	in_data, err := json.Marshal(data)
	gate_data, gatelen := gate.GateMarshal()
	tmpbuf := make([]byte, 4)
	pklen := uint32(len(in_data)) + 4 + gatelen | msghandle.msgHead
	binary.BigEndian.PutUint32(tmpbuf, pklen)
	buff.Write(tmpbuf)
	buff.Write(gate_data)
	buff.Write(in_data)
	return buff.Bytes(), err
}

//编码
func (msghandle *GateJsonMessageHandle) Marshal(msgid uint32, data interface{}) ([]byte, error) {
	return nil, nil
	// buff := &bytes.Buffer{}
	// in_data, err := json.Marshal(data)
	// tmpbuf := make([]byte, 4)
	// pklen := uint32(len(in_data)+8) | msghandle.msgHead
	// binary.BigEndian.PutUint32(tmpbuf, pklen)
	// buff.Write(tmpbuf)
	// binary.BigEndian.PutUint32(tmpbuf, msgid)
	// buff.Write(tmpbuf)
	// buff.Write(in_data)
	// return buff.Bytes(), err
}

//解码
func (msghandle *GateJsonMessageHandle) Unmarshal(buff []byte) (data interface{}, err error) {
	// read := bytes.NewBuffer(buff)

	pklen := binary.BigEndian.Uint32(buff[:4])
	pklen = pklen ^ msghandle.msgHead
	if pklen != uint32(len(buff)) {
		return nil, fmt.Errorf("MsgLen Error:%d.", pklen)
	}
	buff = buff[4:]
	msgid := binary.BigEndian.Uint32(buff[:4])
	msget, err := msghandle.GetRoute(msgid)
	if err != nil {
		return nil, err
	}
	/*
		1、本服务器消息 可以是 IGateMessage
			一般消息结构就可以
		2、客户端过来的消息 IGateChange
			一般消息结构，本地不解，转发给别的服务器
		3、服务器转发的消息 IGateChange,IGateMessage
			二级消息结构，消息头告知转发动作，
			如果是单转，应该可以不用重新打包
			如果是多转，要把消息改成单转进行转发
	*/
	if changemsg, ok := msget.(IGateChange); ok {
		if gatemsg, ok := msget.(IGateMessage); ok {
			//服务器转发消息
			buff, _ := gatemsg.GateUnmarshal(buff)
			changemsg.SetBuffByte(buff)
			return msget, err
		} else if bmsg, ok := msget.(IMessage); ok {
			//客户端消息
			bmsg.SetAction(msgid)
			buff = buff[4:]
			changemsg.SetBuffByte(buff)
			return msget, err
		}
	} else if gatemsg, ok := msget.(IGateMessage); ok {
		//本服务器的消息
		buff, _ := gatemsg.GateUnmarshal(buff)
		err = json.Unmarshal(buff, msget)
		return msget, err
	} else if bmsg, ok := msget.(IMessage); ok {
		//本服务器的消息
		buff = buff[4:]
		err = json.Unmarshal(buff, msget)
		bmsg.SetAction(msgid)
		return msget, err
	}
	return nil, errors.New(fmt.Sprintf("route not IMessage.MsgID:%d", msgid))

}

//设置消息路由
func (msghandle *GateJsonMessageHandle) SetRoute(msgid uint32, msg interface{}) {
	msghandle.routelist[msgid] = msg
}

//按消息拿出消息处理实例
func (msghandle *GateJsonMessageHandle) GetRoute(msgid uint32) (result interface{}, err error) {
	if msget, ok := msghandle.routelist[msgid]; ok {
		return util.ReflectNew(msget)
	} else if msghandle.defType != nil {
		//当没找到的时候，拿出默认的类型
		return util.ReflectNew(msghandle.defType)
	}
	return nil, fmt.Errorf("Not exist MsgID:%d.", msgid)
}

//一个消息是否收完了
//返回这个消息应该的长度，和是否收完的信息
func (msghandle *GateJsonMessageHandle) CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool) {
	pklen := binary.BigEndian.Uint32(buff[:4])
	pklen = pklen ^ msghandle.msgHead
	if pklen > uint32(msghandle.msgMaxLen) {
		return 0, false
	}
	if pklen > uint32(len(buff)) {
		return pklen, false
	}
	return pklen, true
}
