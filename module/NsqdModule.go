package module

import (
	"github.com/buguang01/gsframe/event"
	"github.com/buguang01/Logger"
	"github.com/buguang01/gsframe/threads"
	"github.com/buguang01/gsframe/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/nsqio/go-nsq"
)

type NsqdConfig struct {
	Addr           string //地址
	NSQLookupdAddr string //nsqlookup 地址
	ChanNum        int    //通道缓存空间
}

func NewNsqdModule(configmd *NsqdConfig, sid int) *NsqdModule {
	result := new(NsqdModule)
	result.cg = *configmd
	result.ServerID = util.NewStringInt(sid).ToString()
	return result
}

type NsqdModule struct {
	mgGo     *threads.ThreadGo       //子协程管理器
	getnum   int64                   //收到的总消息数
	sendnum  int64                   //发出去的消息
	consumer *nsq.Consumer           //消费者
	producer *nsq.Producer           //生产者
	chanList chan *event.NsqdMessage //收到的消息
	cg       NsqdConfig              //配置
	ServerID string                  //服务器
	RouteFun event.NsqdHander
}

//Init 初始化
func (this *NsqdModule) Init() {
	this.getnum = 0
	this.sendnum = 0
	this.chanList = make(chan *event.NsqdMessage, this.cg.ChanNum)
	this.mgGo = threads.NewThreadGo()
	this.producer, _ = nsq.NewProducer(this.cg.Addr, nsq.NewConfig())
	this.consumer, _ = nsq.NewConsumer(this.ServerID,
		fmt.Sprintf("%s_channel", this.ServerID), nsq.NewConfig())
	this.consumer.AddHandler(this)
}

//Start 启动
func (this *NsqdModule) Start() {
	this.mgGo.Go(this.Handle)
	if err := this.consumer.ConnectToNSQLookupd(this.cg.NSQLookupdAddr); err != nil {
		panic(err)
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
							msg.SendSID = this.ServerID
							topic := msg.Topic
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
				if msg.Topic == this.ServerID {
					//是发给自己服务器的
					this.mgGo.Go(func(ctx context.Context) {
						this.RouteFun(msg)
					})
				} else {
					//发给别的服务器的
					msg.SendSID = this.ServerID
					topic := msg.Topic
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
		msg := new(event.NsqdMessage)
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
func (this *NsqdModule) AddMsg(msg *event.NsqdMessage) bool {
	msg.SendSID = this.ServerID
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
func (this *NsqdModule) AddMsgSync(msg *event.NsqdMessage) error {
	select {
	case <-this.mgGo.Ctx.Done():
		return errors.New("ctx done")
	default:
		msg.SendSID = this.ServerID
		topic := msg.Topic
		buf, _ := json.Marshal(msg)
		if err := this.producer.Publish(topic, buf); err != nil {
			return err
		}
		atomic.AddInt64(&this.sendnum, 1)

	}
	return nil
}

func (this *NsqdModule) PingNsq(ctx context.Context) bool {
	k := 0
	for {
		if err := this.producer.Ping(); err == nil {
			return true
		} else {
			k = (k + 1) % 10
			Logger.PError(err, "Nsqd Producer Pring Error")
			continue
		}
	}
}
