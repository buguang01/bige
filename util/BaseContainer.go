package util

type BaseContainer struct {
	BaseData
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
