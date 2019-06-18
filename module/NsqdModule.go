package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/buguang01/Logger"
	"github.com/buguang01/gsframe/event"
	"github.com/buguang01/gsframe/threads"
	"github.com/buguang01/util"

	"github.com/nsqio/go-nsq"
)

type NsqdConfig struct {
	Addr                []string //地址
	NSQLookupdAddr      []string //nsqlookup 地址
	ChanNum             int      //通道缓存空间
	LookupdPollInterval int      //去请求lookup nsq节点信息的时间（毫秒）
	MaxInFlight         int      //可以同时访问的节点数

}

func NewNsqdModule(configmd *NsqdConfig, sid int) *NsqdModule {
	result := new(NsqdModule)
	result.cg = *configmd
	result.ServerID = util.NewStringInt(sid).ToString()
	return result
}

type NsqdModule struct {
	mgGo      *threads.ThreadGo       //子协程管理器
	getnum    int64                   //收到的总消息数
	sendnum   int64                   //发出去的消息
	consumer  *nsq.Consumer           //消费者
	producer  *nsq.Producer           //生产者
	chanList  chan event.INsqdMessage //收到的消息
	cg        NsqdConfig              //配置
	ServerID  string                  //服务器
	RouteFun  event.NsqdHander
	GetNewMsg func() event.INsqdMessage //拿到消息接口对象
}

//Init 初始化
func (this *NsqdModule) Init() {
	if this.GetNewMsg == nil {
		this.GetNewMsg = func() event.INsqdMessage {
			return new(event.NsqdMessage)
		}
	}
	this.getnum = 0
	this.sendnum = 0
	this.chanList = make(chan event.INsqdMessage, this.cg.ChanNum)
	this.mgGo = threads.NewThreadGo()
	var err error
	this.producer, err = nsq.NewProducer(this.cg.Addr[0], nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	nsqcg := nsq.NewConfig()
	nsqcg.LookupdPollInterval = time.Duration(this.cg.LookupdPollInterval) * time.Millisecond
	nsqcg.MaxInFlight = this.cg.MaxInFlight
	this.consumer, _ = nsq.NewConsumer(this.ServerID,
		fmt.Sprintf("%s_channel", this.ServerID), nsqcg)
	this.consumer.SetLogger(&nsqlogger{}, nsq.LogLevelError)
	this.consumer.AddHandler(this)

}

//Start 启动
func (this *NsqdModule) Start() {
	this.mgGo.Go(this.Handle)
	if err := this.consumer.ConnectToNSQLookupds(this.cg.NSQLookupdAddr); err != nil {
		panic(err)
	}
	{
		this.registerTopic()
	}
	Logger.PStatus("Nsqd Module Start!")
}

//Stop 停止
func (this *NsqdModule) Stop() {
	this.consumer.Stop()
	<-this.consumer.StopChan
	this.mgGo.CloseWait()
	this.producer.Stop()
	Logger.PStatus("Nsqd Module Stop!")
}

//PrintStatus 打印状态
func (this *NsqdModule) PrintStatus() string {
	return fmt.Sprintf(
		"\r\n		Nsqd Module         :%d/%d/%d	(chanlen/sendnum/getnum)",
		len(this.chanList),
		atomic.AddInt64(&this.sendnum, 0),
		atomic.AddInt64(&this.getnum, 0))
}

//StopConsumer如果要关服，需要提前关闭收消息
func (this *NsqdModule) StopConsumer() {
	this.consumer.Stop()
	<-this.consumer.StopChan
}

func (this *NsqdModule) Handle(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			{
				//要保证所有消息都发出去了。放在消息队列里就可以了
				//因为之前的逻辑可能已做了。
				for {
					select {
					case msg := <-this.chanList:
						{
							msg.SetSendSID(this.ServerID)
							topic := msg.GetTopic()
							buf, _ := json.Marshal(msg)
							if err := this.producer.Publish(topic, buf); err != nil {
								for this.PingNsq(ctx) == true {
									if err := this.producer.Publish(topic, buf); err != nil {
										Logger.PFatal(err)
										continue
									} else {
										break
									}
								}
							}
						}
					default:
						{
							return
						}
					}
				}
			}
		case msg, ok := <-this.chanList:
			{
				if !ok {
					break
				}
				if msg.GetTopic() == this.ServerID {
					//是发给自己服务器的
					this.mgGo.Go(func(ctx context.Context) {
						this.RouteFun(msg)
					})
				} else {
					//发给别的服务器的
					msg.SetSendSID(this.ServerID)
					topic := msg.GetTopic()
					buf, _ := json.Marshal(msg)
					if err := this.producer.Publish(topic, buf); err != nil {
						for this.PingNsq(ctx) == true {
							if err := this.producer.Publish(topic, buf); err != nil {
								Logger.PFatal(err)
								continue
							} else {
								break
							}
						}
					}
					atomic.AddInt64(&this.sendnum, 1)

				}
			}
		}
	}
}

//nsq.Handler的接口
func (this *NsqdModule) HandleMessage(message *nsq.Message) (err error) {
	// fmt.Println(string(message.Body))
	err = nil
	this.mgGo.Try(func(ctx context.Context) {
		if len(message.Body) <= 5 {
			return
		}
		msg := this.GetNewMsg()
		err = json.Unmarshal(message.Body, msg)
		if err != nil {
			Logger.PError(err, "nsqd:%s", string(message.Body))
		} else {
			atomic.AddInt64(&this.getnum, 1)
			this.mgGo.Go(func(ctx context.Context) {
				this.RouteFun(msg)
			})
		}

	}, nil, nil)
	return nil
	//看了源码，如果返回错误，会重新发过来，看nsq的配置
}

//AddMsg 发送消息出去
func (this *NsqdModule) AddMsg(msg event.INsqdMessage) bool {
	msg.SetSendSID(this.ServerID)
	select {
	case <-this.mgGo.Ctx.Done():
		return false
	default:
		this.chanList <- msg
		return true
	}
	// topic := msg.Topic
	// buf, _ := json.Marshal(msg)
	// if err := this.producer.Publish(topic, buf); err != nil {
	// 	//如果出错了就到队列里去
	// }
}

//AddMsgSync 同步发消息出去
func (this *NsqdModule) AddMsgSync(msg event.INsqdMessage) error {
	select {
	case <-this.mgGo.Ctx.Done():
		return errors.New("ctx done")
	default:
		msg.SetSendSID(this.ServerID)
		topic := msg.GetTopic()
		buf, _ := json.Marshal(msg)
		if err := this.producer.Publish(topic, buf); err != nil {
			return err
		}
		atomic.AddInt64(&this.sendnum, 1)

	}
	return nil
}

func (this *NsqdModule) registerTopic() {
	if err := this.producer.Publish(this.ServerID, []byte(" ")); err != nil {
		panic(err)
	}
}

func (this *NsqdModule) PingNsq(ctx context.Context) bool {
	k := 0
	for {
		if err := this.producer.Ping(); err == nil {
			return true
		} else {
			k = (k + 1) % 10
			if k == 0 {
				Logger.PError(err, "Nsqd Producer Pring Error")
				//要换个连接
				for _, addr := range this.cg.Addr {
					if p, err := nsq.NewProducer(addr, nsq.NewConfig()); err == nil {
						if err = p.Ping(); err == nil {
							this.producer.Stop()
							this.producer = p
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
