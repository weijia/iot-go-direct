package iot_go_test

import (
	"encoding/json"
	"testing"

	"iot_go/pkg/msg"
)

func TestConfigMsg(t *testing.T) {
	params := msg.ConfigParam{
		Project: "test",
	}

	config := &msg.ConfigRequest{
		Method: "config",
		Params: params,
	}
	config_json_bytes, _ := json.Marshal(config)
	msg.HandleMsg(config_json_bytes)
}
