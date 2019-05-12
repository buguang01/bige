package util

//BaseBinary DB上的2进制数据
type BaseBinary struct {
	Data          []uint8 //放数据的数组
	ArrayLen      int     //上面那个数组的长度
	OneDataBitNum uint8   //一个数据占多少个位，只能是1，2，4，8
}

func NewBinary(alen int) *BaseBinary {
	result := new(BaseBinary)
	result.Data = make([]uint8, alen)
	result.ArrayLen = alen
	result.OneDataBitNum = 1
	return result
}
func NewBinaryByLen(alen, dlen int) *BaseBinary {
	result := new(BaseBinary)

	result.Data = make([]uint8, alen)
	result.ArrayLen = alen
	result.OneDataBitNum = uint8(dlen)
	return result
}

func (this *BaseBinary) Init(d []uint8) {
	this.Data = d
	this.ArrayLen = len(this.Data)
}

func (this *BaseBinary) ContainKey(index int) int {
	dlen := int(this.OneDataBitNum)
	ai := index / (8 / dlen)
	bitnum := uint8(index % (8 / dlen) * dlen)
	if ai >= this.ArrayLen {
		return 0
	}

	d := this.Data[ai]
	if this.OneDataBitNum == 8 {
		return int(d)
	}
	var mask uint8 = (1 << this.OneDataBitNum) - 1
	return int((d & (mask << bitnum)) >> bitnum)
}

func (this *BaseBinary) UpData(index, val int) bool {
	bval := uint8(val)

	dlen := int(this.OneDataBitNum)
	ai := (index / (8 / dlen))
	bitnum := uint8(index % (8 / dlen) * dlen)
	if ai >= this.ArrayLen {
		return false
	}
	key := this.Data[ai]

	if this.OneDataBitNum == 8 {
		key = bval
	} else {
		var tmp uint8 = (1 << this.OneDataBitNum) - 1
		bval = bval & tmp
		tmp = tmp << bitnum
		tmp = key & ^tmp

		key = uint8(tmp + (bval << bitnum))
	}
	this.Data[ai] = key
	return true

}

func (this *BaseBinary) ToValuesJson() []interface{} {
	dlen := int(this.OneDataBitNum)
	result := make([]interface{}, this.ArrayLen*8/dlen)
	index := 0
	for k := 0; k < this.ArrayLen; k++ {
		key := this.Data[k]
		if this.OneDataBitNum == 8 {
			result[index] = key
			index++
		} else {

			for bitnum := uint8(0); bitnum < 8; bitnum = bitnum + this.OneDataBitNum {
				var mask uint8 = (1 << this.OneDataBitNum) - 1
				result[index] = (key & (mask << bitnum)) >> bitnum
				index++
			}
		}

	}
	return result
}
