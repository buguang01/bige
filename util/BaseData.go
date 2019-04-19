package util

import tsgutils "github.com/typa01/go-utils"

//BaseData 基础仓库数据
type BaseData map[int64]int64

func (this BaseData) UpData(key, num int64) {
	v, _ := this[key]
	if num+v > 0 {
		this[key] = v + num
	} else {
		delete(this, key)
	}
}

func (this BaseData) ToString() string {
	sb := tsgutils.NewStringBuilder()
	t := 0
	for k, v := range this {
		if t == 0 {
			t++
		} else {
			sb.Append(";")
		}
		sb.AppendInt64(k)
		sb.Append(";")
		sb.AppendInt64(v)
	}
	return sb.ToString()
}
