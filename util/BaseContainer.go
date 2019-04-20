package util

type BaseContainer struct {
	BaseData
}

func NewBaseContainer(str string) *BaseContainer {
	result := new(BaseContainer)
	result.BaseData = NewBaseDataString(str)
	return result
}

func (this *BaseContainer) UpItem(itemid, num, max int) {
	v, _ := this.BaseData[itemid]
	if num+v > 0 {
		this.BaseData[itemid] = v + num
		if this.BaseData[itemid] > max {
			this.BaseData[itemid] = max
		}
	} else {
		delete(this.BaseData, itemid)
	}
}
