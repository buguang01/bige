package modules

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util/threads"
	"github.com/nsqio/go-nsq"
)

type NsqdModule struct {
	nsqdPorts          []string                         //nsqd地址组
	lookupPorts        []string                         //lookup地址组
	chanNum            int                              //发消息出去的缓存大小
	lookupPollInterval time.Duration                    //去请求lookup nsq节点信息的时间（秒）
	maxInFlight        int                              //可以同时访问的nsqd节点数
	ServerID           string                           //服务器
	RouteHandle        messages.IMessageHandle          //消息路由
	sendList           chan messages.INsqdResultMessage //发出去的消息
	getnum             int64                            //收到的总消息数
	sendnum            int64                            //发出去的消息
	tmpnum             int64                            //临时计数
	consumer           *nsq.Consumer                    //消费者
	producer           *nsq.Producer                    //生产者
	thgo               *threads.ThreadGo                //子协程管理
}

func NewNsqdModule(opts ...options) *NsqdModule {
	result := &NsqdModule{
		nsqdPorts:          []string{":4150"},
		lookupPorts:        []string{":4161"},
		chanNum:            1024,
		lookupPollInterval: 1 * time.Second,
		maxInFlight:        1024,
		ServerID:           "0",
		RouteHandle:        messages.JsonMessageHandleNew(),
		getnum:             0,
		sendnum:            0,
		tmpnum:             0,
		thgo:               threads.NewThreadGo(),
	}

	for _, opt := range opts {
		opt(result)
	}
	return result
}

//Init 初始化
func (mod *NsqdModule) Init() {
	mod.sendList = make(chan messages.INsqdResultMessage, mod.chanNum)
	var err error
	mod.producer, err = nsq.NewProducer(mod.nsqdPorts[0], nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	nsqcg := nsq.NewConfig()
	nsqcg.LookupdPollInterval = mod.lookupPollInterval
	nsqcg.MaxInFlight = mod.maxInFlight
	mod.consumer, err = nsq.NewConsumer(mod.ServerID,
		fmt.Sprintf("%s_channel", mod.ServerID), nsqcg)
	if err != nil {
		panic(err)
	}
	mod.consumer.SetLogger(&nsqlogger{}, nsq.LogLevelError)
	mod.consumer.AddHandler(mod)
}

//Start 启动
func (mod *NsqdModule) Start() {
	mod.thgo.Go(mod.Handle)
	if err := mod.consumer.ConnectToNSQLookupds(mod.lookupPorts); err != nil {
		panic(err)
	}
	{
		mod.registerTopic()
	}
	Logger.PStatus("Nsqd Module Start!")
}

//Stop 停止
func (mod *NsqdModule) Stop() {
	mod.consumer.Stop()
	<-mod.consumer.StopChan
	mod.thgo.CloseWait()
	mod.producer.Stop()
	Logger.PStatus("Nsqd Module Stop!")
}

//PrintStatus 打印状态
func (mod *NsqdModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n\t\tNsqd Module\t:%d/%d/%d\t(chanlen/sendnum/getnum)",
		len(mod.sendList),
		atomic.LoadInt64(&mod.sendnum),
		atomic.LoadInt64(&mod.getnum))
}

func (mod *NsqdModule) Handle(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				/*
					当关闭服务的时候，
					如果有协程进入了发消息的逻辑里就先等一下
				*/
				if atomic.LoadInt64(&mod.tmpnum) == 0 {
					return
				}
			}
		case msg := <-mod.sendList:
			{
				atomic.AddInt64(&mod.tmpnum, -1)
				if msg.GetTopic() == mod.ServerID {
					atomic.AddInt64(&mod.sendnum, 1)
					//是发给自己服务器的
					mod.thgo.Try(func(ctx context.Context) {
						modmsg, ok := msg.(messages.INsqMessageHandle)
						if !ok {
							Logger.PInfo("Nsqd Send Self Not is Nsqd Msg:%+v", msg)
							return
						} else {
							Logger.PInfo("Nsqd Send Self Msg:%+v", msg)
						}
						atomic.AddInt64(&mod.getnum, 1)
						mod.thgo.Go(func(ctx context.Context) {
							modmsg.NsqDirectCall()
						})
					}, nil, nil)
				} else {
					//发给别的服务器的
					msg.SetSendSID(mod.ServerID)
					topic := msg.GetTopic()
					buf, _ := mod.RouteHandle.Marshal(msg.GetAction(), msg)
					if err := mod.producer.Publish(topic, buf); err != nil {
						for mod.PingNsq(ctx) == true {
							if err := mod.producer.Publish(topic, buf); err != nil {
								Logger.PFatal(err)
								continue
							} else {
								break
							}
						}
					}
					atomic.AddInt64(&mod.sendnum, 1)

				}
			}
		}
	}
}

//AddMsg 发送消息出去
func (mod *NsqdModule) AddMsg(msg messages.INsqdResultMessage) bool {
	msg.SetSendSID(mod.ServerID)
	atomic.AddInt64(&mod.tmpnum, 1)
	select {
	case <-mod.thgo.Ctx.Done():
		atomic.AddInt64(&mod.tmpnum, -1)
		return false
	default:
		mod.sendList <- msg
		return true
	}
}

//AddMsgSync 同步发消息出去
func (mod *NsqdModule) AddMsgSync(msg messages.INsqdResultMessage) error {
	atomic.AddInt64(&mod.tmpnum, 1)
	defer atomic.AddInt64(&mod.tmpnum, -1)
	select {
	case <-mod.thgo.Ctx.Done():
		return errors.New("ctx done")
	default:
		msg.SetSendSID(mod.ServerID)
		topic := msg.GetTopic()
		buf, _ := mod.RouteHandle.Marshal(msg.GetAction(), msg)

		if err := mod.producer.Publish(topic, buf); err != nil {
			return err
		}
		atomic.AddInt64(&mod.sendnum, 1)

	}
	return nil
}

//nsq.Handler的接口
//收nsqd的消息
func (mod *NsqdModule) HandleMessage(message *nsq.Message) (err error) {
	// fmt.Println(string(message.Body))
	err = nil
	mod.thgo.Try(func(ctx context.Context) {
		if len(message.Body) <= 5 {
			return
		}
		buff := message.Body
		// fmt.Println(string(buff))
		msg, err := mod.RouteHandle.Unmarshal(buff)
		if err != nil {
			Logger.PInfo("Nsqd RouteHandle Unmarshal Error:%s", err.Error())
			return
		}
		modmsg, ok := msg.(messages.INsqMessageHandle)
		if !ok {
			Logger.PInfo("Not is Nsq Msg:%+v", msg)
			return
		} else {
			Logger.PInfo("Nsq Get Msg:%+v", msg)
		}
		atomic.AddInt64(&mod.getnum, 1)
		mod.thgo.Go(func(ctx context.Context) {
			modmsg.NsqDirectCall()
		})
	}, nil, nil)
	return nil
	//看了源码，如果返回错误，会重新发过来，看nsq的配置
}

func (mod *NsqdModule) registerTopic() {
	if err := mod.producer.Publish(mod.ServerID, []byte(" ")); err != nil {
		panic(err)
	}
}

func (mod *NsqdModule) PingNsq(ctx context.Context) bool {
	k := 0
	for {
		if err := mod.producer.Ping(); err == nil {
			return true
		} else {
			k = (k + 1) % 10
			if k == 0 {
				Logger.PError(err, "Nsqd Producer Pring Error")
				//要换个连接
				for _, addr := range mod.nsqdPorts {
					if p, err := nsq.NewProducer(addr, nsq.NewConfig()); err == nil {
						if err = p.Ping(); err == nil {
							mod.producer.Stop()
							mod.producer = p
							break
						}
					}
				}
			}
			time.Sleep(1 * time.Second)

			continue
		}
	}
}

type nsqlogger struct{}

func (this *nsqlogger) Output(calldepth int, s string) error {
	Logger.PrintLog(&Logger.LogMsgModel{
		Msg:   s,
		LogLv: Logger.LogLevelstatuslevel + 1,
		Stack: "",
		KeyID: -1,
	})
	return nil
}
