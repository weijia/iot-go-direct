package iot_go_test

import (
	"fmt"
	"testing"

	"iot_go/pkg/msg"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

func TestNodePacket(t *testing.T) {
	msg.DumpBytes(node.GetUpdateGlassColorMsg(util.DecodeId("01020304"), "122365"))
}

func TestChecksum(t *testing.T) {
	data := node.GetUpdateGlassColorMsg(util.DecodeId("01020304"), "122365")
	msg.DumpBytes(data)
	// asset check result using go
	if node.IsChecksumCorrect(data) {
		fmt.Printf("checksum correct\n")
	} else {
		fmt.Printf("checksum incorrect\n")
	}
}

func TestNodeInit(t *testing.T) {
	m := shared.Module{Freq: 470, Band: 250, Factor: 9}
	data := node.GetNodeInitMsg("01020304", m)
	msg.DumpBytes(data)
	// asset check result using go
	if node.IsChecksumCorrect(data) {
		fmt.Printf("checksum correct\n")
	} else {
		fmt.Printf("checksum incorrect\n")
	}
}
