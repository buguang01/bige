package model

import (
	"github.com/buguang01/gsframe/loglogic"
	"time"

	"github.com/garyburd/redigo/redis"
)

//RedisConfigModel 配置信息
type RedisConfigModel struct {
	ConAddr     string        //连接字符串
	MaxIdle     int           //最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	MaxActive   int           //最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
	IdleTimeout time.Duration //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭(秒)
	Password    string        //连接密码
	// Wait bool  //是否等待，设计中应该都是要等待的，所以就不开放了。
}

//RedisAccess redis 管理器
type RedisAccess struct {
	DBConobj *redis.Pool //redis连接池
	cg       *RedisConfigModel
}

//NewRedisAccess 生成新的管理器
func NewRedisAccess(conf *RedisConfigModel) *RedisAccess {
	result := new(RedisAccess)
	result.cg = conf
	result.DBConobj = redis.NewPool(result.dial, result.cg.MaxIdle)
	result.DBConobj.MaxActive = result.cg.MaxActive
	result.DBConobj.IdleTimeout = result.cg.IdleTimeout * time.Second
	result.DBConobj.Wait = true
	result.DBConobj.TestOnBorrow = result.testOnBorrow
	return result
}

func (access *RedisAccess) dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", access.cg.ConAddr)
	if err != nil {
		return nil, err
	}
	if access.cg.Password == "" {
		loglogic.PDebug("redis dial.")
		return c, err
	}
	if _, err := c.Do("AUTH", access.cg.Password); err != nil {
		c.Close()
		return nil, err
	}
	loglogic.PDebug("redis dial.")
	return c, err
}

func (access *RedisAccess) testOnBorrow(c redis.Conn, t time.Time) error {
	loglogic.PDebug("redis testOnBorrow.")
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
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
	loglogic.PDebug("redis close.")

}

//RedisHandleModel 自己把reids的一些常用命令写在这里
type RedisHandleModel struct {
	redis.Conn
}

//Set 写入指定的KEY，val，还有时间dt；如果dt==-1，表示没有时间
func (rd *RedisHandleModel) Set(key string, val interface{}, dt int64) (reply interface{}, err error) {
	if dt > 0 {
		return rd.Do("set", key, val, "EX", dt)
	}
	return rd.Do("set", key, val)

}

//Get 读指定key的值
func (rd *RedisHandleModel) Get(key string) (reply interface{}, err error) {
	return rd.Do("get", key)
}

//DictSet 写入指定(字典\map)表中的指定的KEY，val，还有时间dt；如果dt==-1，表示没有时间
func (rd *RedisHandleModel) DictSet(dname, key string, val interface{}, dt int64) (reply interface{}, err error) {
	if dt > 0 {
		return rd.Do("hset", dname, key, val, "EX", dt)
	}
	return rd.Do("hset", dname, key, val)

}

//DictGet 读指定(字典\map)表中的指定key的值
func (rd *RedisHandleModel) DictGet(dname, key string) (reply interface{}, err error) {
	return rd.Do("hget", dname, key)
}
