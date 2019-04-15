package event

import (
	"buguang01/gsframe/loglogic"
	"encoding/json"
	"net/http"
)

//HTTPcall http的调用方法定义
type HTTPcall func(et JsonMap, w http.ResponseWriter)

//HTTPReplyMsg 回复消息
func HTTPReplyMsg(w http.ResponseWriter, et JsonMap, resultcom int32, jsdata interface{}) {
	jsresult := make(JsonMap)
	jsresult["ACTION"] = et.GetAction()
	// jsresult["ACTIONKEY"] = et["ACTIONKEY"]
	jsresult["ACTIONCOM"] = resultcom
	if jsdata != nil {
		jsresult["JSDATA"] = jsdata
	} else {
		jsresult["JSDATA"] = struct{}{}
	}
	b, _ := json.Marshal(jsresult)
	loglogic.PInfo(string(b))
	w.Write(b)
}
