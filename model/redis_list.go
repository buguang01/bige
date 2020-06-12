package model

import (
	"github.com/garyburd/redigo/redis"
)

/*
LIST就是一个链表，有左头和右头，左右都可以写入数据，也都可以读出数据
通过下面的命令可以对其操作形成，先进先出，先进后出，等队列
*/

/*
LPUSH key value [value ...]
将一个或多个值 value 插入到列表 key 的表头
如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表头：
比如说，对空列表 mylist 执行命令 LPUSH mylist a b c ，列表的值将是 c b a ，
这等同于原子性地执行 LPUSH mylist a 、 LPUSH mylist b 和 LPUSH mylist c 三个命令。
如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。
返回值：
执行 LPUSH 命令后，列表的长度。
*/
func (rd *RedisHandleModel) Lpush(key, val string) (int, error) {
	return redis.Int(rd.Do("LPUSH", key, val))
}

/*
LPUSH key value [value ...]
将一个或多个值 value 插入到列表 key 的表头
如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表头：
比如说，对空列表 mylist 执行命令 LPUSH mylist a b c ，列表的值将是 c b a ，
这等同于原子性地执行 LPUSH mylist a 、 LPUSH mylist b 和 LPUSH mylist c 三个命令。
如果 key 不存在，一个空列表会被创建并执行 LPUSH 操作。
返回值：
执行 LPUSH 命令后，列表的长度。
*/
func (rd *RedisHandleModel) Lpushs(key string, vals ...string) (int, error) {
	p := []interface{}{
		key,
	}
	for _, v := range vals {
		p = append(p, v)
	}
	return redis.Int(rd.Do("LPUSH", p...))
}

/*
LPUSHX key value
将值 value 插入到列表 key 的表头，当且仅当 key 存在并且是一个列表。
和 LPUSH 命令相反，当 key 不存在时， LPUSHX 命令什么也不做。
返回值：
LPUSHX 命令执行之后，表的长度。
*/
func (rd *RedisHandleModel) LpushX(key, val string) (int, error) {
	return redis.Int(rd.Do("LPUSHX", key, val))
}

/*
LPOP key
移除并返回列表 key 的头元素。
返回值：
列表的头元素。
当 key 不存在时，返回 nil 。
*/
func (rd *RedisHandleModel) Lpop(key string) (string, error) {
	return redis.String(rd.Do("LPOP", key))
}

/*
LINDEX key index
返回列表 key 中，下标为 index 的元素。
下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
如果 key 不是列表类型，返回一个错误。
*/
func (rd *RedisHandleModel) Lindex(key string, index int) (string, error) {
	return redis.String(rd.Do("LINDEX", key, index))
}

/*
LSET key index value
将列表 key 下标为 index 的元素的值设置为 value 。
当 index 参数超出范围，或对一个空列表( key 不存在)进行 LSET 时，返回一个错误。
关于列表下标的更多信息，请参考 LINDEX 命令。
时间复杂度：
对头元素或尾元素进行 LSET 操作，复杂度为 O(1)。
其他情况下，为 O(N)， N 为列表的长度。
返回值：
操作成功返回 ok ，否则返回错误信息。
*/
func (rd *RedisHandleModel) Lset(key string, index int, val string) (string, error) {
	return redis.String(rd.Do("LSET", key, index, val))
}

/*
LLEN key
返回列表 key 的长度。
如果 key 不存在，则 key 被解释为一个空列表，返回 0 .
如果 key 不是列表类型，返回一个错误。
*/
func (rd *RedisHandleModel) LLen(key string) (int, error) {
	return redis.Int(rd.Do("LLEN", key))
}

/*
LINSERT key BEFORE|AFTER pivot value
将值 value 插入到列表 key 当中，位于值 pivot 之前或之后。
当 pivot 不存在于列表 key 时，不执行任何操作。
当 key 不存在时， key 被视为空列表，不执行任何操作。
如果 key 不是列表类型，返回一个错误。
时间复杂度:
O(N)， N 为寻找 pivot 过程中经过的元素数量。
返回值:
如果命令执行成功，返回插入操作完成之后，列表的长度。
如果没有找到 pivot ，返回 -1 。
如果 key 不存在或为空列表，返回 0 。
*/
func (rd *RedisHandleModel) Linsert(key string, site RedisListSite, pivot, val string) (int, error) {
	return redis.Int(rd.Do("LINSERT", key, site, pivot, val))
}

/*
LRANGE key start stop
返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
例子：
执行 LRANGE list 0 10 ，结果是一个包含11个元素的列表
时间复杂度:
O(S+N)， S 为偏移量 start ， N 为指定区间内元素的数量。
返回值:
一个列表，包含指定区间内的元素。
*/
func (rd *RedisHandleModel) Lrange(key string, start, stop int) ([]string, error) {
	return redis.Strings(rd.Do("LRANGE", key, start, stop))
}

/*
LREM key count value
根据参数 count 的值，移除列表中与参数 value 相等的元素。
count 的值可以是以下几种：
count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。
count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
count = 0 : 移除表中所有与 value 相等的值。
时间复杂度：
O(N)， N 为列表的长度。
返回值：
被移除元素的数量。
因为不存在的 key 被视作空表(empty list)，所以当 key 不存在时， LREM 命令总是返回 0 。
*/
func (rd *RedisHandleModel) Lrem(key string, count int, val string) (int, error) {
	return redis.Int(rd.Do("LREM", key, count, val))
}

/*
LTRIM key start stop
对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
当 key 不是列表类型时，返回一个错误。
超出范围的下标
超出范围的下标值不会引起错误。
如果 start 下标比列表的最大下标 end ( LLEN list 减去 1 )还要大，或者 start > stop ， LTRIM 返回一个空列表(因为 LTRIM 已经将整个列表清空)。
如果 stop 下标比 end 下标还要大，Redis将 stop 的值设置为 end 。
时间复杂度:
O(N)， N 为被移除的元素的数量。
返回值:
命令执行成功时，返回 ok 。
*/
func (rd *RedisHandleModel) Ltrim(key string, start, stop int) (string, error) {
	return redis.String(rd.Do("LTRIM", key, start, stop))
}

/*
RPOP key
移除并返回列表 key 的尾元素。
返回值：
列表的尾元素。
当 key 不存在时，返回 nil 。
*/
func (rd *RedisHandleModel) Rpop(key string) (string, error) {
	return redis.String(rd.Do("RPOP", key))
}

/*
RPOPLPUSH source destination
命令 RPOPLPUSH 在一个原子时间内，执行以下两个动作：
将列表 source 中的最后一个元素(尾元素)弹出，并返回给客户端。
将 source 弹出的元素插入到列表 destination ，作为 destination 列表的的头元素。
举个例子，你有两个列表 source 和 destination ， source 列表有元素 a, b, c ， destination 列表有元素 x, y, z ，
执行 RPOPLPUSH source destination 之后， source 列表包含元素 a, b ， destination 列表包含元素 c, x, y, z ，并且元素 c 会被返回给客户端。
如果 source 不存在，值 nil 被返回，并且不执行其他动作。
如果 source 和 destination 相同，则列表中的表尾元素被移动到表头，并返回该元素，可以把这种特殊情况视作列表的旋转(rotation)操作。
时间复杂度：
O(1)
返回值：
被弹出的元素。
*/
func (rd *RedisHandleModel) RpopLpush(source, dest string) (string, error) {
	return redis.String(rd.Do("RPOPLPUSH", source, dest))
}

/*
RPUSH key value [value ...]
将一个或多个值 value 插入到列表 key 的表尾(最右边)。
如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，
得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。
当 key 存在但不是列表类型时，返回一个错误。
返回值：
执行 RPUSH 操作后，表的长度。
*/
func (rd *RedisHandleModel) Rpush(key, val string) (int, error) {
	return redis.Int(rd.Do("RPUSH", key, val))
}

/*
RPUSH key value [value ...]
将一个或多个值 value 插入到列表 key 的表尾(最右边)。
如果有多个 value 值，那么各个 value 值按从左到右的顺序依次插入到表尾：比如对一个空列表 mylist 执行 RPUSH mylist a b c ，
得出的结果列表为 a b c ，等同于执行命令 RPUSH mylist a 、 RPUSH mylist b 、 RPUSH mylist c 。
如果 key 不存在，一个空列表会被创建并执行 RPUSH 操作。
当 key 存在但不是列表类型时，返回一个错误。
返回值：
执行 RPUSH 操作后，表的长度。
*/
func (rd *RedisHandleModel) Rpushs(key string, vals ...string) (int, error) {
	p := []interface{}{
		key,
	}
	for _, v := range vals {
		p = append(p, v)
	}
	return redis.Int(rd.Do("RPUSH", p...))
}

/*
RPUSHX key value
将值 value 插入到列表 key 的表尾，当且仅当 key 存在并且是一个列表。
和 RPUSH 命令相反，当 key 不存在时， RPUSHX 命令什么也不做。
返回值：
RPUSHX 命令执行之后，表的长度。
*/
func (rd *RedisHandleModel) RpushX(key, val string) (int, error) {
	return redis.Int(rd.Do("RPUSHX", key, val))
}

/*
BLPOP key [key ...] timeout
简答说就是从这几个key的列表中读一个数据出来，如果没有就等待（阻塞）时间为timeout，
timeout==0为无限等待；
返回数据为：从哪个key读出来的，读出来的值，可能的错误；
说明：
BLPOP 是列表的阻塞式(blocking)弹出原语。
它是 LPOP 命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被 BLPOP 命令阻塞，直到等待超时或发现可弹出元素为止。
当给定多个 key 参数时，按参数 key 的先后顺序依次检查各个列表，弹出第一个非空列表的头元素。
例：
当 BLPOP 被调用时，如果给定 key 内至少有一个非空列表，那么弹出遇到的第一个非空列表的头元素，并和被弹出元素所属的列表的名字一起，组成结果返回给调用者。
当存在多个给定 key 时， BLPOP 按给定 key 参数排列的先后顺序，依次检查各个列表。
假设现在有 job 、 command 和 request 三个列表，其中 job 不存在， command 和 request 都持有非空列表。考虑以下命令：
BLPOP job command request 0
BLPOP 保证返回的元素来自 command ，因为它是按”查找 job -> 查找 command -> 查找 request “这样的顺序，第一个找到的非空列表。
*/
func (rd *RedisHandleModel) Blpop(timeout string, keys ...string) (string, string, error) {
	p := make([]interface{}, len(keys)+1)
	for i, v := range keys {
		p[i] = v
	}
	p[len(keys)] = timeout
	result, err := redis.Strings(rd.Do("BLPOP", p...))
	return result[0], result[1], err
}

/*
BRPOP key [key ...] timeout
BRPOP 是列表的阻塞式(blocking)弹出原语。
与BLPOP一样，只是弹出元素的位置不同
*/
func (rd *RedisHandleModel) Brpop(timeout string, keys ...string) (string, string, error) {
	p := make([]interface{}, len(keys)+1)
	for i, v := range keys {
		p[i] = v
	}
	p[len(keys)] = timeout
	result, err := redis.Strings(rd.Do("BRPOP", p...))
	return result[0], result[1], err
}
