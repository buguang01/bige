package event

import (
	"buguang01/gsframe/loglogic"
	"buguang01/gsframe/threads"
	"encoding/json"

	"golang.org/x/net/websocket"
)

//WebSocketModel 用户连接对象
type WebSocketModel struct {
	*websocket.Conn
	CloseFun WebSocketClose //关闭连接时的方法
	ConInfo  interface{}    //自定义的连接信息，给上层逻辑使用
}

//WebSocketCall websocket调用方法定义
type WebSocketCall func(et JsonMap, wsmd *WebSocketModel, runobj *threads.ThreadGo)

//WebSocketClose 用户连接关闭时的方法
type WebSocketClose func(wsmd *WebSocketModel)

//WebSocketReplyMsg 回复消息
func WebSocketReplyMsg(wsmd *WebSocketModel, et JsonMap, resultcom int32, jsdata JsonMap) {
	jsresult := make(JsonMap)
	jsresult["ACTION"] = et["ACTION"]
	// jsresult["ACTIONKEY"] = et["ACTIONKEY"]
	jsresult["ACTIONCOM"] = resultcom
	if jsdata != nil {
		jsresult["JSDATA"] = jsdata
	} else {
		jsresult["JSDATA"] = struct{}{}
	}
	b, _ := json.Marshal(jsresult)
	loglogic.PInfo(string(b))
	wsmd.Write(b)
}

//WebSocketSendMsg 主动给一个用户发消息
func WebSocketSendMsg(wsmd *WebSocketModel, action int32, jsdata JsonMap) {
	jsresult := make(JsonMap)
	jsresult["ACTION"] = action
	// jsresult["ACTIONKEY"] = et["ACTIONKEY"]
	jsresult["ACTIONCOM"] = 0
	if jsdata != nil {
		jsresult["JSDATA"] = jsdata
	} else {
		jsresult["JSDATA"] = struct{}{}
	}
	b, _ := json.Marshal(jsresult)
	loglogic.PInfo(string(b))
	wsmd.Write(b)
}
