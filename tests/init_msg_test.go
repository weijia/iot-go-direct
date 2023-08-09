package iot_go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"iot_go/pkg/msg"
)

func TestInitMsg(t *testing.T) {

	init := &msg.Init{
		MsgType:       "init",
		GatewayNodeID: "testing",
	}
	init_json_bytes, _ := json.Marshal(init)
	// msg.HandleMsg(init_json_bytes)
	fmt.Println(init_json_bytes)
}
