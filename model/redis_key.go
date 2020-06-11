package model

import "github.com/garyburd/redigo/redis"

//MIGRATE 		不实现
//MOVE
//RANDOMKEY
//RESTORE
//SORT

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

/*
EXPIREAT key timestamp
EXPIREAT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
不同在于 EXPIREAT 命令接受的时间参数是 UNIX 时间戳(unix timestamp)。

例：
redis> EXPIREAT cache 1355292000     # 这个 key 将在 2012.12.12 过期
(integer) 1
redis> TTL cache
(integer) 45081860
*/
func (rd *RedisHandleModel) Expireat(key string, unix int) (int, error) {
	return redis.Int(rd.Do("EXPIREAT", key, unix))
}

/*
KEYS pattern

查找所有符合给定模式 pattern 的 key 。

KEYS * 匹配数据库中所有 key 。
KEYS h?llo 匹配 hello ， hallo 和 hxllo 等。
KEYS h*llo 匹配 hllo 和 heeeeello 等。
KEYS h[ae]llo 匹配 hello 和 hallo ，但不匹配 hillo 。
特殊符号用 \ 隔开
注：KEYS 的速度非常快，但在一个大的数据库中使用它仍然可能造成性能问题，
如果你需要从一个数据集中查找特定的 key ，你最好还是用 Redis 的集合结构(set)来代替。
返回值：
符合给定模式的 key 列表。
*/
func (rd *RedisHandleModel) Keys(pat string) ([]string, error) {
	return redis.Strings(rd.Do("KEYS", pat))
}

/*
PERSIST key
移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。
返回值：
当生存时间移除成功时，返回 1 .
如果 key 不存在或 key 没有设置生存时间，返回 0 。
*/
func (rd *RedisHandleModel) Persist(key string) (int, error) {
	return redis.Int(rd.Do("PERSIST", key))
}

/*
RENAME key newkey

将 key 改名为 newkey 。
当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。
当 newkey 已经存在时， RENAME 命令将覆盖旧值。
返回值：
改名成功时提示 OK ，失败时候返回一个错误。
*/
func (rd *RedisHandleModel) Rename(key, newkey string) (string, error) {
	return redis.String(rd.Do("RENAME", key, newkey))
}

/*
RENAMENX key newkey

当且仅当 newkey 不存在时，将 key 改名为 newkey 。
当 key 不存在时，返回一个错误。
返回值：
修改成功时，返回 1 。
如果 newkey 已经存在，返回 0 。
*/
func (rd *RedisHandleModel) RenameNx(key, newkey string) (int, error) {
	return redis.Int(rd.Do("RENAMENX", key, newkey))
}

/*
TTL key
以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
返回值：
当 key 不存在时，返回 -2 。
当 key 存在但没有设置剩余生存时间时，返回 -1 。
否则，以秒为单位，返回 key 的剩余生存时间。
*/
func (rd *RedisHandleModel) TTL(key string) (int, error) {
	return redis.Int(rd.Do("TTL", key))
}

/*
TYPE key
返回 key 所储存的值的类型。
返回值：
none (key不存在)
string (字符串)
list (列表)
set (集合)
zset (有序集)
hash (哈希表)
*/
func (rd *RedisHandleModel) Type(key string) (string, error) {
	return redis.String(rd.Do("TYPE", key))

}

/*
SCAN 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素（a collection of elements）：

SCAN 命令用于迭代当前数据库中的数据库键。
ZSCAN 命令用于迭代有序集合中的元素（包括元素成员和元素分值）。
*/

/*
SCAN cursor [MATCH pattern] [COUNT count]

SCAN 命令是一个基于游标的迭代器（cursor based iterator）：
SCAN 命令每次被调用之后， 都会向用户返回一个新的游标，
用户在下次迭代时需要使用这个新游标作为 SCAN 命令的游标参数， 以此来延续之前的迭代过程。
当 SCAN 命令的游标参数被设置为 0 时， 服务器将开始一次新的迭代，
而当服务器向用户返回值为 0 的游标时， 表示迭代已结束。

返回：
1) "17"
2)  1) "key:12"
    2) "key:8"
    3) "key:4"
    4) "key:14"
    5) "key:16"
    6) "key:17"
    7) "key:15"
    8) "key:10"
    9) "key:3"
    10) "key:7"
    11) "key:1"
*/
func (rd *RedisHandleModel) Scan(cursor int, match string, count int) ([]string, error) {
	p := []interface{}{
		cursor,
	}
	if match != "" {
		p = append(p, match)
	}
	//count默认值为10
	if count != 10 {
		p = append(p, count)
	}
	return redis.Strings(rd.Do("SCAN", p...))
}

/*
SSCAN key cursor [MATCH pattern] [COUNT count]
命令用于迭代集合键中的元素。
*/
func (rd *RedisHandleModel) SScan(key string, cursor int, match string, count int) ([]string, error) {
	p := []interface{}{
		key, cursor,
	}
	if match != "" {
		p = append(p, match)
	}
	//count默认值为10
	if count != 10 {
		p = append(p, count)
	}
	return redis.Strings(rd.Do("SSCAN", p...))
}

/*
HSCAN key cursor [MATCH pattern] [COUNT count]
命令用于迭代哈希键中的键值对。
*/
func (rd *RedisHandleModel) HScan(key string, cursor int, match string, count int) ([]string, error) {
	p := []interface{}{
		key, cursor,
	}
	if match != "" {
		p = append(p, match)
	}
	//count默认值为10
	if count != 10 {
		p = append(p, count)
	}
	return redis.Strings(rd.Do("HSCAN", p...))
}

/*
ZSCAN key cursor [MATCH pattern] [COUNT count]
命令用于迭代有序集合中的元素（包括元素成员和元素分值）。
*/
func (rd *RedisHandleModel) ZScan(key string, cursor int, match string, count int) ([]string, error) {
	p := []interface{}{
		key, cursor,
	}
	if match != "" {
		p = append(p, match)
	}
	//count默认值为10
	if count != 10 {
		p = append(p, count)
	}
	return redis.Strings(rd.Do("ZSCAN", p...))
}
