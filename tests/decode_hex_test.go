package iot_go_test

import (
	"fmt"
	"iot_go/pkg/util"
	"testing"
)

func TestDecodeHex(t *testing.T) {
	s := "12"
	i := util.GetGlassAreaFromStr(s[0])
	fmt.Printf("Result: %d\n", i)
}
