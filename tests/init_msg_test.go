package iot_go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"iot_go/pkg/shared"
)

func TestInitMsg(t *testing.T) {

	init := shared.Init{
		MsgType: "init",
		InitMsgContent: shared.InitMsgContent{
			NodeInfoContent: shared.NodeInfoContent{
				GatewayNodeID: "testing",
			},
		},
	}
	init_json_bytes, _ := json.Marshal(&init)
	// msg.HandleMsg(init_json_bytes)
	fmt.Println(init_json_bytes)
}
