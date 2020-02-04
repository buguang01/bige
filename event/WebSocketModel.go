package event

import (
	"github.com/buguang01/Logger"
	"encoding/json"
	"github.com/buguang01/util/threads"

	"golang.org/x/net/websocket"
)

//WebSocketModel 用户连接对象
type WebSocketModel struct {
	*websocket.Conn
	CloseFun WebSocketClose //关闭连接时的方法
	ConInfo  interface{}    //自定义的连接信息，给上层逻辑使用
	KeyID    int            //用来标记的ID
}

//WebSocketCall websocket调用方法定义
type WebSocketCall func(et JsonMap, wsmd *WebSocketModel, runobj *threads.ThreadGo)

//WebSocketClose 用户连接关闭时的方法
type WebSocketClose func(wsmd *WebSocketModel)

//WebSocketReplyMsg 回复消息
func WebSocketReplyMsg(wsmd *WebSocketModel, et JsonMap, resultcom int, jsdata JsonMap) {
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
	Logger.PInfoKey(string(b), wsmd.KeyID)
	wsmd.Write(b)
}

//WebSocketSendMsg 主动给一个用户发消息
func WebSocketSendMsg(wsmd *WebSocketModel, action int, jsdata JsonMap) {
	if wsmd == nil {
		Logger.PDebug("WebSocket is nil.Action:%d,Jsdata:%v.", action, jsdata)
		return
	}
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
	Logger.PInfoKey(string(b), wsmd.KeyID)
	wsmd.Write(b)
}
