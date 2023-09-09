package iot_go_test

import (
	"fmt"
	"testing"

	"iot_go/pkg/node"
)

func DumpBytes(a []byte) {
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d, %d, %02x %c\n", i, a[i], a[i], a[i])
	}
	fmt.Printf("\n")
}

func TestNodePacket(t *testing.T) {
	DumpBytes(node.GetUpdateGlassColorMsg("Nod1", "122365"))
}

func TestChecksum(t *testing.T) {
	data := node.GetUpdateGlassColorMsg("Nod1", "122365")
	DumpBytes(data)
	// asset check result using go
	if node.IsChecksumCorrect(data) {
		fmt.Printf("checksum correct\n")
	} else {
		fmt.Printf("checksum incorrect\n")
	}
}
