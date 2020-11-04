package model

import "github.com/garyburd/redigo/redis"

/*
SET key value [EX seconds] [PX milliseconds] [NX|XX]
EX second ：设置键的过期时间为 second 秒。 SET key value EX second 效果等同于 SETEX key second value 。EX==-1表示不设置过期时间
PX millisecond ：设置键的过期时间为 millisecond 毫秒。 SET key value PX millisecond 效果等同于 PSETEX key millisecond value 。(不使用)
NX ：只在键不存在时，才对键进行设置操作。 SET key value NX 效果等同于 SETNX key value 。
XX ：只在键已经存在时，才对键进行设置操作。
*/
func (rd *RedisHandleModel) Set(key string, val string, ex int64, n RedisSetParam) (interface{}, error) {
	parli := make([]interface{}, 0, 5)
	parli = append(parli, key, val)
	if ex != -1 {
		parli = append(parli, "EX", ex)
	}
	switch n {
	case Set_NX:
		parli = append(parli, "NX")
	case Set_XX:
		parli = append(parli, "XX")
	}
	return rd.Do("SET", parli...)
}

/*
当 key 不存在时，返回 nil ，否则，返回 key 的值。
如果 key 不是字符串类型，那么返回一个错误。
*/
func (rd *RedisHandleModel) Get(key string) (string, error) {
	return redis.String(rd.Do("GET", key))
}

/*
STRLEN key
返回 key 所储存的字符串值的长度。
当 key 储存的不是字符串值时，返回一个错误。
返回值：
字符串值的长度。
当 key 不存在时，返回 0 。
*/
func (rd *RedisHandleModel) StrLen(key string) (int, error) {
	return redis.Int(rd.Do("STRLEN", key))
}

/*
如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
返回值：
追加 value 之后， key 中字符串的长度。
*/
func (rd *RedisHandleModel) Append(key, val string) (int, error) {
	return redis.Int(rd.Do("APPEND", key, val))
}

/*
GETRANGE key start end
返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。
GETRANGE 通过保证子字符串的值域(range)不超过实际字符串的值域来处理超出范围的值域请求。
时间复杂度：
O(N)， N 为要返回的字符串的长度。
复杂度最终由字符串的返回值长度决定，但因为从已有字符串中取出子字符串的操作非常廉价(cheap)，所以对于长度不大的字符串，该操作的复杂度也可看作O(1)。
返回值：
截取得出的子字符串。
*/
func (rd *RedisHandleModel) GetRange(key string, start, end int) (string, error) {
	return redis.String(rd.Do("GETRANGE", key, start, end))
}

/*
SETRANGE key offset value
用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
不存在的 key 当作空白字符串处理。
SETRANGE 命令会确保字符串足够长以便将 value 设置在指定的偏移量上，如果给定 key 原来储存的字符串长度比偏移量小(比如字符串只有 5 个字符长，
	但你设置的 offset 是 10 )，那么原字符和偏移量之间的空白将用零字节(zerobytes, "\x00" )来填充。
注意你能使用的最大偏移量是 2^29-1(536870911) ，因为 Redis 字符串的大小被限制在 512 兆(megabytes)以内。如果你需要使用比这更大的空间，你可以使用多个 key 。
当生成一个很长的字符串时，Redis 需要分配内存空间，该操作有时候可能会造成服务器阻塞(block)。
在2010年的Macbook Pro上，设置偏移量为 536870911(512MB 内存分配)，耗费约 300 毫秒，
设置偏移量为 134217728(128MB 内存分配)，耗费约 80 毫秒，设置偏移量 33554432(32MB 内存分配)，
耗费约 30 毫秒，设置偏移量为 8388608(8MB 内存分配)，耗费约 8 毫秒。
注意若首次内存分配成功之后，再对同一个 key 调用 SETRANGE 操作，无须再重新内存。
时间复杂度：
对小(small)的字符串，平摊复杂度O(1)。(关于什么字符串是”小”的，请参考 APPEND 命令)
否则为O(M)， M 为 value 参数的长度。
返回值：
被 SETRANGE 修改之后，字符串的长度。

例：
# 对空字符串/不存在的 key 进行 SETRANGE
redis> EXISTS empty_string
(integer) 0
redis> SETRANGE empty_string 5 "Redis!"   # 对不存在的 key 使用 SETRANGE
(integer) 11
redis> GET empty_string                   # 空白处被"\x00"填充
"\x00\x00\x00\x00\x00Redis!"
*/
func (rd *RedisHandleModel) SetRange(key string, offset int, val string) (int, error) {
	return redis.Int(rd.Do("SETRANGE", key, offset, val))

}

/*
GETSET key value
将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
当 key 存在但不是字符串类型时，返回一个错误。
返回值：
返回给定 key 的旧值。
当 key 没有旧值时，也即是， key 不存在时，返回 nil 。
*/
func (rd *RedisHandleModel) GetSet(key, val string) (string, error) {
	return redis.String(rd.Do("GETSET", key, val))
}

/*
返回所有(一个或多个)给定 key 的值。
如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
返回值：
一个包含所有给定 key 的值的列表。
*/
func (rd *RedisHandleModel) MGet(keys ...interface{}) ([]string, error) {
	return redis.Strings(rd.Do("MGET", keys...))
}

/*
MSET key value [key value ...]
同时设置一个或多个 key-value 对。
返回值：
总是返回 OK (因为 MSET 不可能失败)
*/
func (rd *RedisHandleModel) MSet(keyvals ...interface{}) {
	rd.Do("MSET", keyvals)
}

/*
MSETNX key value [key value ...]
同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一
性逻辑对象(unique logic object)，所有字段要么全被设置，要么全不被设置。
返回值：
当所有 key 都成功设置，返回 1 。
如果所有给定 key 都设置失败(至少有一个 key 已经存在)，那么返回 0 。
*/
func (rd *RedisHandleModel) MsetNx(keyvals ...interface{}) (int, error) {
	return redis.Int(rd.Do("MSETNX", keyvals))

}

/*
对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
当 key 不存在时，自动生成一个新的字符串值。
字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
返回值：
指定偏移量原来储存的位数据
*/
func (rd *RedisHandleModel) SetBit(key string, offset, val int) (int, error) {
	return redis.Int(rd.Do("SETBIT", key, offset, val))
}

/*
对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
*/
func (rd *RedisHandleModel) GetBit(key string, offset int) (int, error) {
	return redis.Int(rd.Do("GETBIT", key, offset))
}

/*
BITCOUNT key [start] [end]
计算给定字符串中，被设置为 1 的比特位的数量。
一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值：比如 -1 表示最后一个位，而 -2 表示倒数第二个位，以此类推。
不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
返回值：
被设置为 1 的位的数量。

如果 start==end 时，表示所有数据
*/
func (rd *RedisHandleModel) BitCount(key string, start, end int) (int, error) {
	if end == 0 {
		return redis.Int(rd.Do("BITCOUNT", key))
	}
	return redis.Int(rd.Do("BITCOUNT", key, start, end))
}

/*
BITOP operation destkey key [key ...]
对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
operation 可以是 AND 、 OR 、 NOT 、 XOR 这四种操作中的任意一种：
BITOP AND destkey key [key ...] ，对一个或多个 key 求逻辑并，并将结果保存到 destkey 。
BITOP OR destkey key [key ...] ，对一个或多个 key 求逻辑或，并将结果保存到 destkey 。
BITOP XOR destkey key [key ...] ，对一个或多个 key 求逻辑异或，并将结果保存到 destkey 。
BITOP NOT destkey key ，对给定 key 求逻辑非，并将结果保存到 destkey 。
除了 NOT 操作之外，其他操作都可以接受一个或多个 key 作为输入。
处理不同长度的字符串
当 BITOP 处理不同长度的字符串时，较短的那个字符串所缺少的部分会被看作 0 。
空的 key 也被看作是包含 0 的字符串序列。

时间复杂度：
O(N)
返回值：
保存到 destkey 的字符串的长度，和输入 key 中最长的字符串长度相等。
BITOP 的复杂度为 O(N) ，当处理大型矩阵(matrix)或者进行大数据量的统计时，最好将任务指派到附属节点(slave)进行，避免阻塞主节点。
*/

func (rd *RedisHandleModel) Bitop(bitopt RedisBitOperation, deskey string, keys ...string) (int, error) {

	switch bitopt {
	case Bit_And:
		p := []interface{}{
			"AND", deskey,
		}
		for _, v := range keys {
			p = append(p, v)
		}
		return redis.Int(rd.Do("BITOP", p...))
	case Bit_Or:
		p := []interface{}{
			"OR", deskey,
		}
		for _, v := range keys {
			p = append(p, v)
		}
		return redis.Int(rd.Do("BITOP", p...))
	case Bit_Xor:
		p := []interface{}{
			"XOR", deskey,
		}
		for _, v := range keys {
			p = append(p, v)
		}
		return redis.Int(rd.Do("BITOP", p...))
	case Bit_Not:
		if len(keys) == 1 {
			return redis.Int(rd.Do("BITOP", "NOT", deskey, keys[0]))
		}
	}
	return 0, ErrParam
}

/*
将 key 中储存的数字值减一。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
返回值：
执行 DECR 命令之后 key 的值。
*/
func (rd *RedisHandleModel) Decr(key string) (int, error) {
	return redis.Int(rd.Do("DECR", key))
}

/*
INCR key
将 key 中储存的数字值增一。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
返回值：
执行 INCR 命令之后 key 的值。
*/
func (rd *RedisHandleModel) Incr(key string) (int, error) {
	return redis.Int(rd.Do("INCR", key))

}

/*
将 key 所储存的值减去减量 decrement 。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
返回值：
减去 decrement 之后， key 的值。
*/
func (rd *RedisHandleModel) Decrby(key string, dec int) (int, error) {
	return redis.Int(rd.Do("DECRBY", key, dec))
}

/*
INCRBY key increment
将 key 所储存的值加上增量 increment 。
如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
返回值：
加上 increment 之后， key 的值。
*/
func (rd *RedisHandleModel) Incrby(key string, inc int) (int, error) {
	return redis.Int(rd.Do("INCRBY", key, inc))
}

/*
INCRBYFLOAT key increment
为 key 中所储存的值加上浮点数增量 increment 。
如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者。
无论是 key 的值，还是增量 increment ，都可以使用像 2.0e7 、 3e5 、 90e-2 那样的指数符号(exponential notation)来表示，
但是，执行 INCRBYFLOAT 命令之后的值总是以同样的形式储存，也即是，它们总是由一个数字，一个（可选的）小数点和一个任意位的小数部分组成（比如 3.14 、 69.768 ，诸如此类)，
小数部分尾随的 0 会被移除，如果有需要的话，还会将浮点数改为整数（比如 3.0 会被保存成 3 ）。
除此之外，无论加法计算所得的浮点数的实际精度有多长， INCRBYFLOAT 的计算结果也最多只能表示小数点的后十七位。
当以下任意一个条件发生时，返回一个错误：
key 的值不是字符串类型(因为 Redis 中的数字和浮点数都以字符串的形式保存，所以它们都属于字符串类型）
key 当前的值或者给定的增量 increment 不能解释(parse)为双精度浮点数(double precision floating point number）
*/
func (rd *RedisHandleModel) IncrbyFloat(key string, inc float64) (float64, error) {
	return redis.Float64(rd.Do("INCRBYFLOAT", key, inc))
}
