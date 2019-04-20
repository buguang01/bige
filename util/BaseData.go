package util

//BaseData 基础仓库数据
type BaseData map[int]int

func (this BaseData) Clone() BaseData {
	result := make(BaseData)
	for k, v := range this {
		result[k] = v
	}
	return result
}

func NewBaseDataString(str string) BaseData {
	result := make(BaseData)
	if str == "" || str == "0" {
		return result
	}
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

func (this BaseData) UpDataBc(addbc, delbc BaseData) {
	if delbc != nil {
		for k, n := range delbc {
			this.UpData(k, -n)
		}
	}
	if addbc != nil {
		for k, n := range addbc {
			this.UpData(k, n)
		}
	}
}

func (this BaseData) GetNumByKey(key int) int {
	v, ok := this[key]
	if !ok {
		return 0
	}
	return v
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
