package json

import jsoniter "github.com/json-iterator/go"

var JsonLib = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	a, b := JsonLib.Marshal(v)
	return a, b
}

func Unmarshal(data []byte, v interface{}) error {
	return JsonLib.Unmarshal(data, v)
}
