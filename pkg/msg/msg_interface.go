package msg

import (
	"encoding/json"

	"iot_go/pkg/util"
)

type BaseMsg struct {
	Method string `json:"method"`
}

type MsgHandler interface {
	handle()
}

type MsgFactory interface {
	getTargetMsg() any
}

func HandleMsg(body []byte) {
	var base_msg BaseMsg
	if err := json.Unmarshal(body, &base_msg); err != nil {
		util.IotLog(err)
	}
	switch base_msg.Method {
	case "config":
		var config Config
		if err := json.Unmarshal(body, &config); err != nil {
			util.IotLog(err)
		}
		config.handle()
	}
}
