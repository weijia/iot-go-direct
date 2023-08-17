package iot_go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"iot_go/pkg/shared"
)

func TestJsonUnmarshal(t *testing.T) {

	init := shared.Init{
		MsgType: "init",
		InitMsgContent: shared.InitMsgContent{
			NodeInfoContent: shared.NodeInfoContent{
				GatewayNodeID: "testing",
			},
		},
	}
	init_json_str, _ := json.Marshal(&init)
	fmt.Println(string(init_json_str))
	b := []byte(init_json_str)
	var dat shared.Init
	if err := json.Unmarshal(b, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	// fmt.Println("hello test")
}
