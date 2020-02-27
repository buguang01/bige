package messages

import (
	"encoding/json"
	"fmt"

	"github.com/buguang01/util"
)

type HttpJsonMessageHandle struct {
	routelist map[uint32]interface{} //消息路由列表
}

func HttpJsonMessageHandleNew(opts ...options) (msghandle *HttpJsonMessageHandle) {
	msghandle = &HttpJsonMessageHandle{
		routelist: make(map[uint32]interface{}),
	}
	for _, opt := range opts {
		opt(msghandle)
	}
	return msghandle
}

//编码
func (msghandle *HttpJsonMessageHandle) Marshal(msgid uint32, data interface{}) ([]byte, error) {
	in_data, err := json.Marshal(data)
	return in_data, err
}

//解码
func (msghandle *HttpJsonMessageHandle) Unmarshal(buff []byte) (data interface{}, err error) {
	js := make(JsonMap)
	if err = json.Unmarshal(buff, &js); err != nil {
		return nil, err
	}
	data, err = msghandle.GetRoute(js.GetAction())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buff, data)
	return data, err
}

//设置消息路由
func (msghandle *HttpJsonMessageHandle) SetRoute(msgid uint32, msg interface{}) {
	msghandle.routelist[msgid] = msg
}

//按消息拿出消息处理实例
func (msghandle *HttpJsonMessageHandle) GetRoute(msgid uint32) (result interface{}, err error) {
	if msget, ok := msghandle.routelist[msgid]; ok {
		return util.ReflectNew(msget)
	}
	return nil, fmt.Errorf("Not exist MsgID:%d.", msgid)
}

//一个消息是否收完了
//返回这个消息应该的长度，和是否收完的信息
func (msghandle *HttpJsonMessageHandle) CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool) {
	return uint32(len(buff)), true
}
