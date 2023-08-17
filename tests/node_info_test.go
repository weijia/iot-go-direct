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
		GatewayNodeID: "testing",
	}
	json_bytes, _ := json.Marshal(&node_info_request)
	msg.HandleMsg(json_bytes)
	fmt.Println(json_bytes)
}
