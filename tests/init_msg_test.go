package iot_go_test

import (
	"testing"

	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
)

func TestInitMsg(t *testing.T) {

	init := shared.Init{
		MsgType: "init",
		InitMsgContent: shared.InitMsgContent{
			NodeInfoContent: shared.NodeInfoContent{
				GatewayNodeId: "testing",
			},
		},
	}
	msg.DumpMsg(init)
}
