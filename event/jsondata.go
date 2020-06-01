package event

import "github.com/buguang01/util"

type JsonMap map[string]interface{}

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

func (js JsonMap) GetString(key string) *util.String {
	return util.NewStringAny(js[key])
}

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

func (js JsonArray) GetString(index int) *util.String {
	return util.NewStringAny(js[index])
}
