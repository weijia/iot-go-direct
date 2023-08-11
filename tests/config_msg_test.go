package iot_go_test

import (
	"encoding/json"
	"testing"

	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
)

func TestConfigMsg(t *testing.T) {
	params := shared.ConfigParams{
		Project: "test",
	}

	config := &msg.ConfigRequest{
		Method: "config",
		Params: params,
	}
	config_json_bytes, _ := json.Marshal(config)
	msg.HandleMsg(config_json_bytes)
}
