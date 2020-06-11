package model

import (
	"runtime"
	"time"

	"github.com/buguang01/Logger"

	"github.com/garyburd/redigo/redis"
)

var (
	//用来设置默认协程数
	moduleCap = runtime.NumCPU() * 10
)

//设置redis地址,带端口
func RedisSetAddr(addr string) redisoptions {
	return func(mod *RedisAccess) {
		mod.addr = addr
	}
}

//设置redis,auth(进入密码)
func RedisSetAuth(auth string) redisoptions {
	return func(mod *RedisAccess) {
		mod.auth = auth
	}
}

//设置redis默认的indexdb
func RedisSetIndexDB(indexdb int) redisoptions {
	return func(mod *RedisAccess) {
		mod.indexdb = indexdb
	}
}

//设置redis最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
func RedisSetMaxIdle(maxidle int) redisoptions {
	return func(mod *RedisAccess) {
		mod.maxIdle = maxidle
	}
}

//最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
func RedisSetMaxActive(maxactive int) redisoptions {
	return func(mod *RedisAccess) {
		mod.maxActive = maxactive
	}
}

//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭(秒)
func RedisSetIdleTimeout(idleTimeout time.Duration) redisoptions {
	return func(mod *RedisAccess) {
		mod.idleTimeout = idleTimeout * time.Second
	}
}

type redisoptions func(mod *RedisAccess)

//RedisAccess redis 管理器
type RedisAccess struct {
	DBConobj    *redis.Pool   //redis连接池
	addr        string        //连接字符串
	indexdb     int           //默认DB号
	auth        string        //连接密码
	maxIdle     int           //最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	maxActive   int           //最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
	idleTimeout time.Duration //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭(秒)
}

//NewRedisAccess 生成新的管理器
func NewRedisAccess(opts ...redisoptions) *RedisAccess {
	result := &RedisAccess{
		addr:        "127.0.0.1:6379",
		auth:        "",
		maxIdle:     moduleCap / 10,
		maxActive:   moduleCap,
		idleTimeout: 3600 * time.Second,
	}
	for _, opt := range opts {
		opt(result)
	}
	result.DBConobj = redis.NewPool(result.dial, result.maxIdle)
	result.DBConobj.MaxActive = result.maxActive
	result.DBConobj.IdleTimeout = result.idleTimeout
	result.DBConobj.Wait = true
	result.DBConobj.TestOnBorrow = result.testOnBorrow

	return result
}

func (access *RedisAccess) dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", access.addr)
	if err != nil {
		return nil, err
	}
	if access.auth == "" {
		Logger.PDebug("redis dial.")
		return c, err
	}
	if _, err := c.Do("AUTH", access.auth); err != nil {
		c.Close()
		return nil, err
	}
	Logger.PDebug("redis dial.")
	return c, err
}

func (access *RedisAccess) testOnBorrow(c redis.Conn, t time.Time) error {
	// Logger.PDebug("redis testOnBorrow.")
	// if time.Since(t) < time.Minute {
	// 	return nil
	// }
	_, err := c.Do("SELECT", access.indexdb)
	return err
}

//GetConn 拿到一个可用的连接，你要在这句之后写上：defer conn.Close()
//用来在使用完之后把连接放回池子里去
func (access *RedisAccess) GetConn() *RedisHandleModel {
	return &RedisHandleModel{access.DBConobj.Get()}
}

//Close 关闭池子，一般只有关服的时候才用到
func (access *RedisAccess) Close() {
	access.DBConobj.Close()
	Logger.PDebug("redis close.")

}

//RedisHandleModel 自己把reids的一些常用命令写在这里
type RedisHandleModel struct {
	redis.Conn
}

func (rd *RedisHandleModel) Dispose() {
	rd.Close()
}

// //Set 写入指定的KEY，val，还有时间dt；如果dt==-1，表示没有时间
// func (rd *RedisHandleModel) Set(key string, val interface{}, dt int64) (reply interface{}, err error) {
// 	if dt > 0 {
// 		return rd.Do("set", key, val, "EX", dt)
// 	}
// 	return rd.Do("set", key, val)

// }

// //Get 读指定key的值
// func (rd *RedisHandleModel) Get(key string) (reply interface{}, err error) {
// 	return rd.Do("get", key)
// }

// //DictSet 写入指定(字典\map)表中的指定的KEY，val
// func (rd *RedisHandleModel) DictSet(dname, key string, val interface{}) (reply interface{}, err error) {
// 	return rd.Do("hset", dname, key, val)
// }

// //DictGet 读指定(字典\map)表中的指定key的值
// func (rd *RedisHandleModel) DictGet(dname, key string) (reply interface{}, err error) {
// 	return rd.Do("hget", dname, key)
// }

// //DictGetAll 读指定字典表中的所有Key和值
// func (rd *RedisHandleModel) DictGetAll(dname string) (reply interface{}, err error) {
// 	return rd.Do("hgetall", dname)
// }

// //DictGetAllByStringArray 读指定字典表中的所有Key和值,返回[]string
// func (rd *RedisHandleModel) DictGetAllByStringArray(dname string) (result []string, err error) {

// 	reply, err := rd.DictGetAll(dname)
// 	if err != nil {
// 		return nil, err
// 	}
// 	arr, ok := reply.([]interface{})
// 	if !ok {
// 		return nil, errors.New("interface to []interface error.")
// 	}
// 	result = make([]string, len(arr))
// 	for i, v := range arr {
// 		result[i] = string(v.([]byte))
// 	}
// 	return result, nil

// }

// //DelKey 删指定的KEY
// func (rd *RedisHandleModel) DelKey(dname string) (reply interface{}, err error) {
// 	return rd.Do("del", dname)
// }

// //GetKeyList 返回指定KEY的列表，一般用来删除过期的KEY
// func (rd *RedisHandleModel) GetKeyList(dname string) (reply interface{}, err error) {
// 	return rd.Do("scan", 0, "match", dname, "count", 10000)
// }

// //RankGet 写入排行榜
// func (rd *RedisHandleModel) RankSet(rankName, key string, val float64) (reply interface{}, err error) {
// 	return rd.Do("zadd", rankName, val, key)
// }

// //RankAddSet 累加写入排行榜数据
// func (rd *RedisHandleModel) RankAddSet(rankName, key string, val float64) (reply interface{}, err error) {
// 	return rd.Do("zincrby", rankName, val, key)
// }

// //RankGetPage 排行榜多少到多少 从小到大
// func (rd *RedisHandleModel) RankGetPage(rankName string, page1, page2 int) ([]*RankInfoModel, error) {
// 	reply, err := rd.Do("zrange", rankName, page1, page2, "withscores")
// 	if err != nil {
// 		return nil, err
// 	}
// 	arr := reply.([]interface{})
// 	result := make([]*RankInfoModel, len(arr)/2)
// 	rno := page1
// 	for i := 0; i < len(arr); i += 2 {
// 		rno++
// 		md := new(RankInfoModel)
// 		md.RankNo = rno
// 		md.KeyID = util.NewStringAny(arr[i])
// 		md.Val = util.NewStringAny(arr[i+1]).ToFloatV()
// 		result[i/2] = md
// 	}
// 	return result, nil
// }

// //RankRevGetPage 排行榜多少到多少 从大到小
// func (rd *RedisHandleModel) RankRevGetPage(rankName string, page1, page2 int) ([]*RankInfoModel, error) {
// 	reply, err := rd.Do("zrevrange", rankName, page1, page2, "withscores")
// 	if err != nil {
// 		return nil, err
// 	}
// 	arr := reply.([]interface{})
// 	result := make([]*RankInfoModel, len(arr)/2)
// 	rno := page1
// 	for i := 0; i < len(arr); i += 2 {
// 		rno++
// 		md := new(RankInfoModel)
// 		md.RankNo = rno
// 		md.KeyID = util.NewStringAny(arr[i])
// 		md.Val = util.NewStringAny(arr[i+1]).ToFloatV()
// 		result[i/2] = md
// 	}
// 	return result, nil
// }

// //RankDelByKey 删除指定排行榜中的指定key
// func (rd *RedisHandleModel) RankDelByKey(rankName, key string) (reply interface{}, err error) {
// 	return rd.Do("zrem", rankName, key)
// }

// //获取排名（从大到小）
// func (rd *RedisHandleModel) RankRevGetNo(rankName, key string) (result int, err error) {
// 	reply, err := rd.Do("zrevrank", rankName, key)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return util.NewStringAny(reply).ToIntV(), err
// }

// //获取排名（从小到大）
// func (rd *RedisHandleModel) RankGetNo(rankName, key string) (result int, err error) {
// 	reply, err := rd.Do("zrank", rankName, key)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return util.NewStringAny(reply).ToIntV(), err
// }

// //获取指定排行榜上Key的Val
// func (rd *RedisHandleModel) RankGetVal(rankName, key string) (result float64, err error) {
// 	reply, err := rd.Do("zscore", rankName, key)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return util.NewStringAny(reply).ToFloatV(), err
// }

// //获取排行榜从小到大的排名信息
// func (rd *RedisHandleModel) RankGetInfo(rankName, key string) (result *RankInfoModel, err error) {
// 	result = new(RankInfoModel)
// 	result.RankNo, err = rd.RankGetNo(rankName, key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result.Val, err = rd.RankGetVal(rankName, key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// //获取排行榜从大到小的排名信息
// func (rd *RedisHandleModel) RankRevGetInfo(rankName, key string) (result *RankInfoModel, err error) {
// 	result = new(RankInfoModel)
// 	result.RankNo, err = rd.RankRevGetNo(rankName, key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result.Val, err = rd.RankGetVal(rankName, key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// //List的在尾后添加成员
// //"a","b"
// //写入"c"后为：”a“，”b“，”c“；"a"是索引0的位置
// func (rd *RedisHandleModel) ListRpush(rankName, val string) (reply interface{}, err error) {
// 	return rd.Do("rpush", rankName, val)
// }

// //List的在尾后添加成员
// //"a","b"
// //写入"c"后为：”c“,”a“，”b“；"c"是索引0的位置
// func (rd *RedisHandleModel) ListLpush(rankName, val string) (reply interface{}, err error) {
// 	return rd.Do("lpush", rankName, val)
// }

// //List 返回指定索引范围的数据
// func (rd *RedisHandleModel) ListGetArr(rankName string, stindex, num int) ([]string, error) {
// 	reply, err := rd.Do("lrange", rankName, stindex, num)
// 	if err != nil {
// 		return nil, err
// 	}
// 	arr := reply.([]interface{})
// 	result := make([]string, len(arr))
// 	for i, v := range arr {
// 		result[i] = string(v.([]byte))
// 	}
// 	return result, nil
// }

// //List 成员数量
// func (rd *RedisHandleModel) ListLen(rankName string) (int, error) {
// 	reply, err := rd.Do("llen", rankName)
// 	if err != nil {
// 		return -1, err
// 	}
// 	result := util.NewStringAny(reply).ToIntV()
// 	return result, nil
// }

// type RankInfoModel struct {
// 	RankNo int          //名次
// 	KeyID  *util.String //名字
// 	Val    float64      //值
// }
