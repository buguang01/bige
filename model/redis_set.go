package model

import "github.com/garyburd/redigo/redis"

/*
SADD key member [member ...]
将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
当 key 不是集合类型时，返回一个错误。
时间复杂度:
O(N)， N 是被添加的元素的数量。
返回值:
被添加到集合中的新元素的数量，不包括被忽略的元素。
*/
func (rd *RedisHandleModel) Sadd(key string, member string) (int, error) {
	return redis.Int(rd.Do("SADD", key, member))
}

//同Sadd
func (rd *RedisHandleModel) Sadds(key string, members ...string) (int, error) {

	p := []interface{}{
		key,
	}
	for _, m := range members {
		p = append(p, m)
	}
	return redis.Int(rd.Do("SADD", p...))
}

/*
SCARD key
返回集合 key 的基数(集合中元素的数量)。
返回值：
集合的基数。
当 key 不存在时，返回 0 。
*/
func (rd *RedisHandleModel) Scard(key string) (int, error) {
	return redis.Int(rd.Do("SCARD", key))
}

/*
SDIFF key [key ...]
返回一个集合的全部成员，该集合是所有给定集合之间的差集。
不存在的 key 被视为空集。
时间复杂度:
O(N)， N 是所有给定集合的成员数量之和。
返回值:
交集成员的列表。
*/
func (rd *RedisHandleModel) Sdiff(keys ...interface{}) ([]string, error) {
	return redis.Strings(rd.Do("SDIFF", keys...))
}

/*
SDIFFSTORE destination key [key ...]
这个命令的作用和 SDIFF 类似，但它将结果保存到 destination 集合，而不是简单地返回结果集。
如果 destination 集合已经存在，则将其覆盖。
destination 可以是 key 本身。
时间复杂度:
O(N)， N 是所有给定集合的成员数量之和。
返回值:
结果集中的元素数量。
*/
func (rd *RedisHandleModel) SdiffStore(dest string, keys ...string) (int, error) {
	p := []interface{}{
		dest,
	}
	for _, k := range keys {
		p = append(p, k)
	}
	return redis.Int(rd.Do("SDIFFSTORE", p...))
}

/*
SINTER key [key ...]
返回一个集合的全部成员，该集合是所有给定集合的交集。
不存在的 key 被视为空集。
当给定集合当中有一个空集时，结果也为空集(根据集合运算定律)。
时间复杂度:
O(N * M)， N 为给定集合当中基数最小的集合， M 为给定集合的个数。
返回值:
交集成员的列表。
*/
func (rd *RedisHandleModel) Sinter(keys ...interface{}) ([]string, error) {
	return redis.Strings(rd.Do("SINTER", keys...))
}

/*
SINTERSTORE destination key [key ...]
这个命令类似于 SINTER 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
如果 destination 集合已经存在，则将其覆盖。
destination 可以是 key 本身。
时间复杂度:
O(N * M)， N 为给定集合当中基数最小的集合， M 为给定集合的个数。
返回值:
结果集中的成员数量。
*/
func (rd *RedisHandleModel) SinterStore(dest string, keys ...string) (int, error) {
	p := []interface{}{
		dest,
	}
	for _, k := range keys {
		p = append(p, k)
	}
	return redis.Int(rd.Do("SINTERSTORE", p...))
}

/*
SISMEMBER key member
判断 member 元素是否集合 key 的成员。
返回值:
如果 member 元素是集合的成员，返回 1 。
如果 member 元素不是集合的成员，或 key 不存在，返回 0 。
*/
func (rd *RedisHandleModel) Sismember(key, member string) (int, error) {
	return redis.Int(rd.Do("SISMEMBER", key, member))
}

/*
SMEMBERS key
返回集合 key 中的所有成员。
不存在的 key 被视为空集合。
时间复杂度:
O(N)， N 为集合的基数。
返回值:
集合中的所有成员。
*/
func (rd *RedisHandleModel) Smembers(key string) ([]string, error) {
	return redis.Strings(rd.Do("SMEMBERS", key))
}

/*
SMOVE source destination member
将 member 元素从 source 集合移动到 destination 集合。
SMOVE 是原子性操作。
如果 source 集合不存在或不包含指定的 member 元素，则 SMOVE 命令不执行任何操作，仅返回 0 。
否则， member 元素从 source 集合中被移除，并添加到 destination 集合中去。
当 destination 集合已经包含 member 元素时， SMOVE 命令只是简单地将 source 集合中的 member 元素删除。
当 source 或 destination 不是集合类型时，返回一个错误。
返回值:
如果 member 元素被成功移除，返回 1 。
如果 member 元素不是 source 集合的成员，并且没有任何操作对 destination 集合执行，那么返回 0 。
*/
func (rd *RedisHandleModel) Smove(source, dest, member string) (int, error) {
	return redis.Int(rd.Do("SMOVE", source, dest, member))
}

/*
SPOP key
移除并返回集合中的一个随机元素。
如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。
返回值:
被移除的随机元素。
当 key 不存在或 key 是空集时，返回 nil 。
*/
func (rd *RedisHandleModel) Spop(key string) (string, error) {
	return redis.String(rd.Do("SPOP", key))
}

/*
SRANDMEMBER key [count]
如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
从 Redis 2.6 版本开始， SRANDMEMBER 命令接受可选的 count 参数：
如果 count 为正数，且小于集合基数，那么命令返回一个包含 count 个元素的数组，数组中的元素各不相同。如果 count 大于等于集合基数，那么返回整个集合。
如果 count 为负数，那么命令返回一个数组，数组中的元素可能会重复出现多次，而数组的长度为 count 的绝对值。
时间复杂度:
只提供 key 参数时为 O(1) 。
如果提供了 count 参数，那么为 O(N) ，N 为返回数组的元素个数。
返回值:
只提供 key 参数时，返回一个元素；如果集合为空，返回 nil 。
如果提供了 count 参数，那么返回一个数组；如果集合为空，返回空数组。
*/
func (rd *RedisHandleModel) SrandMember(key string, count int) ([]string, error) {
	return redis.Strings(rd.Do("SRANDMEMBER", key, count))
}

/*
SREM key member [member ...]
移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
当 key 不是集合类型，返回一个错误。
返回值:
被成功移除的元素的数量，不包括被忽略的元素。
*/
func (rd *RedisHandleModel) Srem(key string, members ...string) (int, error) {
	p := []interface{}{
		key,
	}
	for _, k := range members {
		p = append(p, k)
	}
	return redis.Int(rd.Do("SREM", p...))
}

/*
SUNION key [key ...]
返回一个集合的全部成员，该集合是所有给定集合的并集。
不存在的 key 被视为空集。
时间复杂度:
O(N)， N 是所有给定集合的成员数量之和。
返回值:
并集成员的列表。
*/
func (rd *RedisHandleModel) Sunion(keys ...interface{}) ([]string, error) {
	return redis.Strings(rd.Do("SUNION", keys))
}

/*
SUNIONSTORE destination key [key ...]
这个命令类似于 SUNION 命令，但它将结果保存到 destination 集合，而不是简单地返回结果集。
如果 destination 已经存在，则将其覆盖。
destination 可以是 key 本身。
时间复杂度:
O(N)， N 是所有给定集合的成员数量之和。
返回值:
结果集中的元素数量。
*/
func (rd *RedisHandleModel) SunionStore(dest string, keys ...string) (int, error) {
	p := []interface{}{
		dest,
	}
	for _, k := range keys {
		p = append(p, k)
	}
	return redis.Int(rd.Do("SUNIONSTORE", p...))
}
