package util

//BaseData 基础仓库数据
type BaseData struct {
	Data map[int]int
}

//Clone 复制一个BaseData
func (this BaseData) Clone() *BaseData {
	result := new(BaseData)
	result.Data = make(map[int]int)
	for k, v := range this.Data {
		result.Data[k] = v
	}
	return result
}

//NewBaseDataString 用字符串初始化一个数据
func NewBaseDataString(str string) *BaseData {
	result := new(BaseData)
	result.Data = make(map[int]int)
	if str == "" || str == "0" {
		return result
	}
	arr := StringToIntArray(str, ";")
	for i := 0; i < len(arr); i += 2 {
		result.Data[arr[i]] = arr[i+1]
	}
	return result
}

//UpData 更新指定数据
func (this *BaseData) UpData(key, num int) {
	v, _ := this.Data[key]
	if num+v > 0 {
		this.Data[key] = v + num
	} else {
		delete(this.Data, key)
	}
}

//UpDataBc批量用别的数据，更新本数据
func (this *BaseData) UpDataBc(addbc, delbc *BaseData) {
	if delbc != nil {
		for k, n := range delbc.Data {
			this.UpData(k, -n)
		}
	}
	if addbc != nil {
		for k, n := range addbc.Data {
			this.UpData(k, n)
		}
	}
}

//GetNumByKey指定数据的值
func (this *BaseData) GetNumByKey(key int) int {
	v, ok := this.Data[key]
	if !ok {
		return 0
	}
	return v
}

//ToString 字符串化
func (this *BaseData) ToString() string {
	sb := NewStringBuilder()
	t := 0
	for k, v := range this.Data {
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

//Count 总数量
func (this *BaseData) Count() (result int) {
	for _, n := range this.Data {
		result += n
	}
	return result
}

//MaxItem 最大数值的KEY，value
func (this *BaseData) MaxItem() (key, num int) {
	for k, n := range this.Data {
		if n > num {
			key, num = k, n
		}
	}
	return key, num
}
