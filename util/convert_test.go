package util_test

import (
	"buguang01/gsframe/util"
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	b := []byte("1001")
	var v int32 = 0
	util.Convert.ByteToAll(b, &v)
	fmt.Println(v)
}
