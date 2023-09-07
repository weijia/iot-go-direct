package iot_go_test

import (
	"encoding/json"
	"testing"

	"iot_go/pkg/bsp"
	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

func TestConfigMsg(t *testing.T) {
	bsp.InitConfig()
	params := shared.ConfigParams{
		BaseConfigParams: shared.BaseConfigParams{
			Project: "test",
		},
		Module1: shared.Module{Freq: 30, Band: 50, Factor: 80},
		Module2: shared.Module{Freq: 35, Band: 55, Factor: 85},
	}

	config := &msg.ConfigRequest{
		Method: "config",
		Params: params,
	}
	config_json_bytes, _ := json.Marshal(config)
	var mqtt util.Mqtt
	reply := msg.HandleMsg(&mqtt, config_json_bytes)
	msg.DumpMsg(reply)
}
