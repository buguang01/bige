package messages

import (
	"github.com/buguang01/util"
)

//JsonMap   收到的JSON数据
type JsonMap map[string]interface{}

//GetAction 消息号
func (js JsonMap) GetAction() uint32 {
	v, _ := util.NewStringAny(js["ACTIONID"]).ToInt()
	return uint32(v)
	// return util.Convert.ToInt32(js["ACTION"])
}

//GetActionKey int32 //消息序号
func (js JsonMap) GetActionKey() int {
	v, _ := util.NewStringAny(js["ACTIONKEY"]).ToInt()
	return v
	// return util.Convert.ToInt32(js["ACTIONKEY"])
}

//GetMemberID int32  //用户ID
func (js JsonMap) GetMemberID() int {
	v, _ := util.NewStringAny(js["MEMBERID"]).ToInt()
	return v
	// return util.Convert.ToInt32(js["MEMBERID"])
}

//GetHash string     //发回给用户信息用的钥匙
func (js JsonMap) GetHash() string {
	v := util.NewStringAny(js["HASH"]).ToString()
	return v
	// return uint32(js["HASH"].(float64))
}

//返回一个JsonArray
func (js JsonMap) GetArray(key string) JsonArray {
	v := js[key].([]interface{})
	return JsonArray(v)
}

//返回一个JsonMap
func (js JsonMap) GetMap(key string) JsonMap {
	v := js[key].(map[string]interface{})
	return JsonMap(v)
}

//返回[]int
func (js JsonMap) GetIntArray(key string) []int {
	return js.GetArray(key).GetIntArray()
}

//JsonArray JSON数组
type JsonArray []interface{}

func (js JsonArray) GetArray(index int) JsonArray {
	v := js[index].([]interface{})
	return JsonArray(v)
}

func (js JsonArray) GetMap(index int) JsonMap {
	v := js[index].(map[string]interface{})
	return JsonMap(v)
}

func (js JsonArray) GetIntArray() []int {
	result := make([]int, len(js))
	for i, v := range js {
		result[i] = util.NewStringAny(v).ToIntV()
	}
	return result
}

type MessageJson struct {
	ActionID uint32 `json:"ACTIONID"`
}
