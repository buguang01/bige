package messages

import (
	"net"
	"net/http"
	"runtime"

	"github.com/buguang01/bige/model"
	"golang.org/x/net/websocket"
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

//WebSocketModel 用户连接对象
type WebSocketModel struct {
	*websocket.Conn
	CloseFun func(wsmd *WebSocketModel) //关闭连接时的方法
	ConInfo  interface{}                //自定义的连接信息，给上层逻辑使用
	KeyID    int                        //用来标记的ID
}
type IWebSocketMessageHandle interface {
	//ws的回调
	WebSocketDirectCall(ws *WebSocketModel)
}

//SocketModel 用户连接对象
type SocketModel struct {
	Conn     net.Conn                //连接信息
	CloseFun func(skmd *SocketModel) //关闭连接时的方法
	ConInfo  interface{}             //自定义的连接信息，给上层逻辑使用
	KeyID    int                     //用来标记的ID
}
type ISocketMessageHandle interface {
	ISocketResultMessage
	//socket的回调
	SocketDirectCall(ws *SocketModel)
}

type ISocketResultMessage interface {
	//消息
	IMessage
	//消息来源的用户ID
	GetSendUserID() int
	//消息来源的服务ID
	GetSendSID() string
	//设置来源服务ID
	SetSendSID(sid string)
	//目标服务器ID
	GetTopic() string
}

type INsqMessageHandle interface {
	INsqdResultMessage
	//Nsq的回调
	NsqDirectCall()
}

type INsqdResultMessage interface {
	//消息
	IMessage
	//消息来源的用户ID
	GetSendUserID() int
	//消息来源的服务ID
	GetSendSID() string
	//设置来源服务ID
	SetSendSID(sid string)
	//目标服务器ID
	GetTopic() string
}

//nsqd消息的基础结构
type NsqdMessage struct {
	SendUserID int    //发信息用户ID
	SendSID    string //发信息服务器（回复用的信息）
	ActionID   uint32 //消息号
	Topic      string //目标
}

func (msg *NsqdMessage) GetSendUserID() int {
	return msg.SendUserID
}
func (msg *NsqdMessage) GetSendSID() string {
	return msg.SendSID
}
func (msg *NsqdMessage) SetSendSID(sid string) {
	msg.SendSID = sid
}
func (msg *NsqdMessage) GetActionID() uint32 {
	return msg.ActionID
}
func (msg *NsqdMessage) GetTopic() string {
	return msg.Topic
}

//LogicMessage 逻辑委托
type ILogicMessage interface {
	//所在协程的KEY
	LogicThreadID() int
	//调用方法
	MessageHandle()
}

type LogicMessage struct {
	UserID int
}

func (msg *LogicMessage) LogicThreadID() int {
	//默认按CPU个数的十倍的协程数来分配对应的协程进行处理
	//分配时，按用户ID进行取余
	cpu := runtime.NumCPU() * 10
	return msg.UserID % cpu
}

//DataBase的处理接口
type IDataBaseMessage interface {
	//所在DB协程
	DBThreadID() int
	/*数据表,如果你的表放入时，不是马上保存的，那么后续可以用这个KEY来进行覆盖，
	这样就可以实现多次修改一次保存的功能
	所以这个字段建议是：用户ID+数据表名+数据主键
	*/
	GetDataKey() string

	//调用方法
	SaveDB(conn model.IConnDB) error
}
