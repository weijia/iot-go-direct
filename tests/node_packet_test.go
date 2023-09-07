package iot_go_test

import (
	"fmt"
	"testing"

	"iot_go/pkg/node"
)

func DumpBytes(a []byte) {
	for i := 0; i < len(a); i++ {
		fmt.Printf("%02x %c\n", a[i], a[i])
	}
	fmt.Printf("\n")
}

func TestNodePacket(t *testing.T) {
	DumpBytes(node.GetUpdateGlassColorMsg("Nod1", "121315"))
}
