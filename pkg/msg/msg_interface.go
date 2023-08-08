package msg

import (
	"encoding/json"

	"iot_go/pkg/util"
)

type BaseMsg struct {
	MsgType string `json:"msg_type"`
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
	switch base_msg.MsgType {
	case "init":
		var init Init
		if err := json.Unmarshal(body, &init); err != nil {
			util.IotLog(err)
		}
		init.handle()
	}
}
