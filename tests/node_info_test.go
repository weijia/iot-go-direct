package iot_go_test

import (
	"context"
	"encoding/json"
	"testing"

	"iot_go/pkg/msg"
)

func TestNodeInfoMsg(t *testing.T) {

	node_info_request := msg.NodeInfoRequest{
		Method: "node_info_request",
		Params: msg.GatewayParams{
			GatewayNodeId: "testing",
		},
	}
	json_bytes, _ := json.Marshal(&node_info_request)
	reply := msg.HandleMqttMsg(context.TODO(), nil, json_bytes)
	msg.DumpMsg(reply)
}
