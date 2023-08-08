package msg

import (
	"encoding/json"
)

type BaseMsg struct {
	MsgType string `json:"msg_type"`
}

type MsgHandler interface {
	handle()
}

func handle(body []byte) {
	var init Init
	json.Unmarshal(body, &init)
}
