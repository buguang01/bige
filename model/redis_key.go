package model

import "github.com/garyburd/redigo/redis"

/*
DEL key [key ...]
删除给定的一个或多个 key 。
不存在的 key 会被忽略。
时间复杂度：
O(N)， N 为被删除的 key 的数量。
删除单个字符串类型的 key ，时间复杂度为O(1)。
删除单个列表、集合、有序集合或哈希表类型的 key ，时间复杂度为O(M)， M 为以上数据结构内的元素数量。
返回值：
被删除 key 的数量。
*/
func (rd *RedisHandleModel) Del(keys ...interface{}) (int, error) {
	return redis.Int(rd.Do("DEL", keys...))
}

/*
EXISTS key
检查给定 key 是否存在。
返回值：
若 key 存在，返回 1 ，否则返回 0 。
*/
func (rd *RedisHandleModel) Exists(key string) (int, error) {
	return redis.Int(rd.Do("EXISTS", key))
}

/*
EXPIRE key seconds
为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。
返回值：
设置成功返回 1 。
当 key 不存在或者不能为 key 设置生存时间时(比如在低于 2.1.3 版本的 Redis 中你尝试更新 key 的生存时间)，返回 0 。
*/
func (rd *RedisHandleModel) Expire(key string, ex int) (int, error) {
	return redis.Int(rd.Do("EXPIRE", key, ex))
}
