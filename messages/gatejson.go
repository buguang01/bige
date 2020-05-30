package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/buguang01/util"
)

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
	if gatemsg, ok := msget.(IGateMessage); ok {
		buff, _ := gatemsg.GateUnmarshal(buff)
		err = json.Unmarshal(buff, msget)
		return msget, err
	} else {
		buff = buff[4:]
		err = json.Unmarshal(buff, msget)
		msget.(IMessage).SetAction(msgid)
		return msget, err
	}

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
