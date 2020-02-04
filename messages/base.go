package messages

import (
	"net/http"

	"github.com/buguang01/bige/event"
)

//消息收发接口
type IMessageHandle interface {
	//编码
	Marshal(msgid uint32, data interface{}) ([]byte, error)
	//解码
	Unmarshal(buff []byte) (data interface{}, err error)
	//设置消息路由
	SetRoute(msgid uint32, msg interface{})
	//按消息拿出消息处理实例
	GetRoute(msgid uint32) (result interface{}, err error)
	//一个消息是否收完了
	CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool)
}

type options func(msghandle IMessageHandle)

type IMessage interface {
	GetAction() uint32
}

type IHttpMessageHandle interface {
	//HTTP的回调
	HttpDirectCall(w http.ResponseWriter, req *http.Request)
}

type IWebSocketMessageHandle interface {
	//ws的回调
	WebSocketDirectCall(ws *event.WebSocketModel)
}

type INsqMessageHandle interface {
	//Nsq的回调
	NsqDirectCall()
}
