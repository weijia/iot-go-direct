package iot_go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"iot_go/pkg/msg"
)

func TestNodeInfoMsg(t *testing.T) {

	node_info_request := msg.NodeInfoRequest{
		Method: "node_info_request",
		Params: msg.GatewayParams{
			GatewayNodeID: "testing",
		},
	}
	json_bytes, _ := json.Marshal(&node_info_request)
	reply := msg.HandleMsg(nil, json_bytes)
	res, _ := json.MarshalIndent(reply, "", "  ")
	fmt.Println(string(res))
}
