package msg

import (
	"encoding/json"
	"testing"

	"iot_go/pkg/msg"
)

func TestInitMsg(t *testing.T) {

	init := &msg.Init{
		MsgType:       "init",
		GatewayNodeID: "testing",
	}
	init_json_bytes, _ := json.Marshal(init)
	msg.HandleMsg(init_json_bytes)
}
