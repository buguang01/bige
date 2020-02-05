package modules

import (
	"database/sql"
	"time"

	"github.com/buguang01/bige/event"
	"github.com/buguang01/bige/messages"
	"github.com/buguang01/util/threads"
)

type DataBaseModule struct {
	chanNum   int                              //通道缓存空间
	timeout   time.Duration                    //超时时间
	logicList map[int]*dataBaseThread          //子逻辑列表
	keyList   []int                            //key列表，用来间隔遍历
	chanList  chan []messages.IDataBaseMessage //消息信道
	getNum    int64                            //收到的总消息数
	endNum    int64                            //处理结束数
	thgo      *threads.ThreadGo                //子协程管理器
}

type dataBaseThread struct {
	DBThreadID int                              //协程ID
	upDataList map[string]event.ISqlDataModel   //缓存要更新的数据
	chanList   chan []messages.IDataBaseMessage //收要更新的数据
	Conndb     *sql.DB                          //数据库连接对象
	upTime     time.Time                        //更新时间
}
