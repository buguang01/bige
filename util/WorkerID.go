package util

import (
	"fmt"
	"sync"
	"time"
)

//SnowFlake的结构体
type SnowFlakeIdGenerator struct {
	workerId           int64 //当前的workerId
	workerIdAfterShift int64 //移位后的workerId，可直接跟时间戳、序号取位或操作
	lastMsTimestamp    int64 //上一次用的时间戳
	curSequence        int64 //当前的序号

	timeBitSize     uint8 //时间戳占的位数，默认为41位，最大不超过60位
	workerIdBitSize uint8 //workerId占的位数，默认10，最大不超过60位
	sequenceBitSize uint8 //序号占的位数，默认12，最大不超过60位

	lock       *sync.Mutex //同步用的
	isHaveInit bool        //是否已经初始化了

	maxWorkerId        int64 //workerId的最大值，初始化时计算出来的
	maxSequence        int64 //最后序列号最大值，初始化时计算出来的
	workerIdLeftShift  uint8 //生成的workerId只取最低的几位，这里要左移，给序列号腾位，初始化时计算出来的
	timestampLeftShift uint8 //生成的时间戳左移几位，给workId、序列号腾位，初始化时计算出来的
}

//实例化一个ID生成器
func NewIDGenerator() *SnowFlakeIdGenerator {
	return &SnowFlakeIdGenerator{
		workerId:           0,
		lastMsTimestamp:    0,
		curSequence:        0,
		timeBitSize:        41, //默认的时间戳占的位数
		workerIdBitSize:    10, //默认的workerId占的位数
		sequenceBitSize:    12, //默认的序号占的位数
		maxWorkerId:        0,  //最大的workerId，初始化时计算出来的
		maxSequence:        0,  //最大的序号值，初始化的时计算出来的
		workerIdLeftShift:  0,  //worker id左移位数
		timestampLeftShift: 0,
		lock:               new(sync.Mutex),
		isHaveInit:         false,
	}
}

//设置worker id
func (sfg *SnowFlakeIdGenerator) SetWorkerId(w int64) *SnowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.workerId = w
	return sfg
}

//设置时间戳占的位数
func (sfg *SnowFlakeIdGenerator) SetTimeBitSize(n uint8) *SnowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.timeBitSize = n
	return sfg
}

//设置worker id占的位数
func (sfg *SnowFlakeIdGenerator) SetWorkerIdBitSize(n uint8) *SnowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.workerIdBitSize = n
	return sfg
}

//设置序号占的位数
func (sfg *SnowFlakeIdGenerator) SetSequenceBitSize(n uint8) *SnowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.sequenceBitSize = n
	return sfg
}

//初始化操作
func (sfg *SnowFlakeIdGenerator) Init() (*SnowFlakeIdGenerator, error) {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()

	//如果已经初始化了
	if sfg.isHaveInit {
		return sfg, nil
	}

	if sfg.sequenceBitSize < 1 || sfg.sequenceBitSize > 60 {
		return nil, fmt.Errorf("Init failed:\tinvalid sequence bit size, should (1,60)")
	}
	if sfg.timeBitSize < 1 || sfg.timeBitSize > 60 {
		return nil, fmt.Errorf("Init failed:\tinvalid time bit size, should (1,60)")
	}
	if sfg.workerIdBitSize < 1 || sfg.workerIdBitSize > 60 {
		return nil, fmt.Errorf("Init failed:\tinvalid worker id bit size, should (1,60)")
	}
	if sfg.workerIdBitSize+sfg.sequenceBitSize+sfg.timeBitSize != 63 {
		return nil, fmt.Errorf("Init failed:\tinvalid sum of all bit size, should eq 63")
	}

	//确定移位数
	sfg.workerIdLeftShift = sfg.sequenceBitSize
	sfg.timestampLeftShift = sfg.sequenceBitSize + sfg.workerIdBitSize

	//确定序列号及workerId最大值
	sfg.maxWorkerId = -1 ^ (-1 << sfg.workerIdBitSize)
	sfg.maxSequence = -1 ^ (-1 << sfg.sequenceBitSize)

	//移位之后的workerId，返回结果时可直接跟时间戳、序号取或操作即可
	sfg.workerIdAfterShift = sfg.workerId << sfg.workerIdLeftShift

	//判断当前的workerId是否合法
	if sfg.workerId > sfg.maxWorkerId {
		return nil, fmt.Errorf("Init failed:\tinvalid worker id, should not greater than %d", sfg.maxWorkerId)
	}

	//初始化完毕
	sfg.isHaveInit = true
	sfg.lastMsTimestamp = 0
	sfg.curSequence = 0
	return sfg, nil
}

//生成时间戳，根据bit size设置取高几位
//即，生成的时间戳先右移几位，再左移几位，就保留了最高的指定位数
func (sfg *SnowFlakeIdGenerator) genTs() int64 {
	rawTs := time.Now().UnixNano()
	diff := 64 - sfg.timeBitSize
	ret := (rawTs >> diff) << diff
	return ret
}

//生成下一个时间戳，如果时间戳的位数较小，且序号用完时此处等待的时间会较长
func (sfg *SnowFlakeIdGenerator) genNextTs(last int64) int64 {
	for {
		cur := sfg.genTs()
		if cur > last {
			return cur
		}
	}
}

//生成下一个ID
func (sfg *SnowFlakeIdGenerator) NextId() (int64, error) {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()

	//如果还没有初始化
	if !sfg.isHaveInit {
		return 0, fmt.Errorf("Gen NextId failed:\tplease execute Init() first")
	}

	//先判断当前的时间戳，如果比上一次的还小，说明出问题了
	curTs := sfg.genTs()
	if curTs < sfg.lastMsTimestamp {
		return 0, fmt.Errorf("Gen NextId failed:\tunknown error, the system clock occur some wrong")
	}

	//如果跟上次的时间戳相同，则增加序号
	if curTs == sfg.lastMsTimestamp {
		sfg.curSequence = (sfg.curSequence + 1) & sfg.maxSequence
		//序号又归0即用完了，重新生成时间戳
		if sfg.curSequence == 0 {
			curTs = sfg.genNextTs(sfg.lastMsTimestamp)
		}
	} else {
		//如果两个的时间戳不一样，则归0序号
		sfg.curSequence = 0
	}

	sfg.lastMsTimestamp = curTs

	//将处理好的各个位组装成一个int64型
	curTs = curTs | sfg.workerIdAfterShift | sfg.curSequence
	return curTs, nil
}

//解析生成的ID
func (sfg *SnowFlakeIdGenerator) Parse(id int64) (int64, int64, int64, error) {
	//如果还没有初始化
	if !sfg.isHaveInit {
		return 0, 0, 0, fmt.Errorf("Parse failed:\tplease execute Init() first")
	}

	//先提取时间戳部分
	shift := sfg.sequenceBitSize + sfg.sequenceBitSize
	timestamp := (id & (-1 << shift)) >> shift

	//再提取workerId部分
	shift = sfg.sequenceBitSize
	workerId := (id & (sfg.maxWorkerId << shift)) >> shift

	//序号部分
	sequence := id & sfg.maxSequence

	//解析错误
	if workerId != sfg.workerId || workerId > sfg.maxWorkerId {
		fmt.Printf("workerBitSize=%d\tMaxWorkerId=%d\n", sfg.workerIdBitSize, sfg.maxWorkerId)
		return 0, 0, 0, fmt.Errorf("parse failed：invalid id, originWorkerId=%d\tparseWorkerId=%d\n",
			sfg.workerId, workerId)
	}
	if sequence < 0 || sequence > sfg.maxSequence {
		fmt.Printf("sequesnceBitSize=%d\tMaxSequence=%d\n", sfg.sequenceBitSize, sfg.maxSequence)
		return 0, 0, 0, fmt.Errorf("parse failed：invalid id, parseSequence=%d\n", sequence)
	}

	return timestamp, workerId, sequence, nil
}
