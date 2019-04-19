package util

//BaseData 基础仓库数据
type BaseData map[int]int

func NewBaseDataString(str string) BaseData {
	result := make(BaseData)
	arr := StringToIntArray(str, ";")
	for i := 0; i < len(arr); i += 2 {
		result[arr[i]] = arr[i+1]
	}
	return result
}

func (this BaseData) UpData(key, num int) {
	v, _ := this[key]
	if num+v > 0 {
		this[key] = v + num
	} else {
		delete(this, key)
	}
}

func (this BaseData) ToString() string {
	sb := NewStringBuilder()
	t := 0
	for k, v := range this {
		if t == 0 {
			t++
		} else {
			sb.Append(";")
		}
		sb.AppendInt(k)
		sb.Append(";")
		sb.AppendInt(v)
	}
	return sb.ToString()
}
