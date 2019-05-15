package util_test

import (
	"github.com/buguang01/gsframe/util"
	"fmt"
	"testing"
)

func TestBinary(t *testing.T) {
	bin := util.NewBinaryByLen(4, 1)
	bin.UpData(3, 1)
	bin.UpData(2, 4)
	fmt.Println(bin.ContainKey(2))
	fmt.Println(bin.ContainKey(0))
	fmt.Println(bin.ToValuesJson())
}
