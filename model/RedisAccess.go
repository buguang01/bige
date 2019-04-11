package model

import (
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
		return c, err
	}
	if _, err := c.Do("AUTH", access.cg.Password); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func (access *RedisAccess) testOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}

//GetConn 拿到一个可用的连接，你要在这句之后写上：defer conn.Close()
//用来在使用完之后把连接放回池子里去
func (access *RedisAccess) GetConn() redis.Conn {
	return access.DBConobj.Get()
}

//Close 关闭池子，一般只有关服的时候才用到
func (access *RedisAccess) Close() {
	access.DBConobj.Close()
}
