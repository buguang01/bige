package model

import "github.com/garyburd/redigo/redis"

/*
HDEL key field [field ...]
删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
返回值:
被成功移除的域的数量，不包括被忽略的域。
*/
func (rd *RedisHandleModel) Hdel(key string, fields ...string) (int, error) {
	p := []interface{}{
		key,
	}
	for _, f := range fields {
		p = append(p, f)
	}

	return redis.Int(rd.Do("HDEL", p))
}

/*
HEXISTS key field
查看哈希表 key 中，给定域 field 是否存在。
返回值：
如果哈希表含有给定域，返回 1 。
如果哈希表不含有给定域，或 key 不存在，返回 0 。
*/
func (rd *RedisHandleModel) Hexists(key, field string) (int, error) {
	return redis.Int(rd.Do("HEXISTS", key, field))
}

/*
HGET key field
返回哈希表 key 中给定域 field 的值。
返回值：
给定域的值。
当给定域不存在或是给定 key 不存在时，返回 nil 。
*/
func (rd *RedisHandleModel) Hget(key, field string) (string, error) {
	return redis.String(rd.Do("HGET", key, field))

}

/*
HGETALL key
返回哈希表 key 中，所有的域和值。
在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
返回值：
以列表形式返回哈希表的域和域的值。
若 key 不存在，返回空列表。

返回值为Map，分为三种不同类型的key：string,int,int64类型
*/
func (rd *RedisHandleModel) HgetallByStringMap(key string) (map[string]string, error) {
	return redis.StringMap(rd.Do("HGETALL", key))
}

/*
同 HgetallByStringMap
*/
func (rd *RedisHandleModel) HgetallByIntMap(key string) (map[string]int, error) {
	return redis.IntMap(rd.Do("HGETALL", key))
}

/*
同 HgetallByStringMap
*/
func (rd *RedisHandleModel) HgetallByInt64Map(key string) (map[string]int64, error) {
	return redis.Int64Map(rd.Do("HGETALL", key))
}

/*
HINCRBY key field increment
为哈希表 key 中的域 field 的值加上增量 increment 。
增量也可以为负数，相当于对给定域进行减法操作。
如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
本操作的值被限制在 64 位(bit)有符号数字表示之内。
返回值：
执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。
*/
func (rd *RedisHandleModel) Hincrby(key, field string, inc int) (int, error) {
	return redis.Int(rd.Do("HINCRBY", key, field, inc))
}

/*
同Hincrby 只是返回值类型为float64
*/
func (rd *RedisHandleModel) HincrbyFloat(key, field string, inc float64) (float64, error) {
	return redis.Float64(rd.Do("HINCRBYFLOAT", key, field, inc))
}

/*
HKEYS key
返回哈希表 key 中的所有域。
时间复杂度：
O(N)， N 为哈希表的大小。
返回值：
一个包含哈希表中所有域的表。
当 key 不存在时，返回一个空表。
*/
func (rd *RedisHandleModel) Hkeys(key string) ([]string, error) {
	return redis.Strings(rd.Do("HKEYS", key))
}

/*
HLEN key
返回哈希表 key 中域的数量。
时间复杂度：
O(1)
返回值：
哈希表中域的数量。
当 key 不存在时，返回 0 。
*/
func (rd *RedisHandleModel) Hlen(key string) (int, error) {
	return redis.Int(rd.Do("HLEN", key))
}

/*
HMGET key field [field ...]
返回哈希表 key 中，一个或多个给定域的值。
如果给定的域不存在于哈希表，那么返回一个 nil 值。
因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。

返回一个与传入field对应的map[string]string
如果对应的值为nil的话，这里的表现为“”
*/
func (rd *RedisHandleModel) Hmget(key string, fields ...string) (map[string]string, error) {
	p := []interface{}{
		key,
	}
	for _, f := range fields {
		p = append(p, f)
	}
	valli, err := redis.Strings(rd.Do("HMGET", p...))
	if err != nil {
		return map[string]string{}, err
	}
	result := make(map[string]string)
	for index, fd := range fields {
		result[fd] = valli[index]
	}
	return result, nil
}

/*
HMSET key field value [field value ...]
同时将多个 field-value (域-值)对设置到哈希表 key 中。
此命令会覆盖哈希表中已存在的域。
如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
返回值：
如果命令执行成功，返回 OK 。
当 key 不是哈希表(hash)类型时，返回一个错误。
*/
func (rd *RedisHandleModel) Hmset(key string, fieldvals ...interface{}) (string, error) {
	p := []interface{}{
		key,
	}
	p = append(p, fieldvals...)

	return redis.String(rd.Do("HMSET", p...))
}

/*
HSETNX key field value
将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
若域 field 已经存在，该操作无效。
如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
返回值：
设置成功，返回 1 。
如果给定域已经存在且没有操作被执行，返回 0 。
*/
func (rd *RedisHandleModel) HsetNx(key, field, val string) (int, error) {
	return redis.Int(rd.Do("HSETNX", key, field, val))
}

/*
HVALS key
返回哈希表 key 中所有域的值。
时间复杂度：
O(N)， N 为哈希表的大小。
返回值：
一个包含哈希表中所有值的表。
当 key 不存在时，返回一个空表。
*/
func (rd *RedisHandleModel) Hvals(key string) ([]string, error) {
	return redis.Strings(rd.Do("HVALS", key))
}
