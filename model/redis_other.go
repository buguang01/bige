package model

/*
SELECT index
切换到指定的数据库，数据库索引号 index 用数字值指定，以 0 作为起始索引值。
默认使用 0 号数据库。
返回 error==nil 就是成功的
*/
func (rd *RedisHandleModel) Select(index int) error {
	_, err := rd.Do("SELECT", index)
	return err
}

/*
AUTH password
通过设置配置文件中 requirepass 项的值(使用命令 CONFIG SET requirepass password )，可以使用密码来保护 Redis 服务器。
返回 error==nil 就是成功的
*/
func (rd *RedisHandleModel) Auth(pwd string) error {
	_, err := rd.Do("AUTH", pwd)
	return err
}

/*
使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
返回 error==nil 就是成功的
*/
func (rd *RedisHandleModel) Ping() error {
	_, err := rd.Do("PING")
	return err
}
