package msg

import (
	"iot_go/pkg/util"
)

type LogConfig struct {
	MsgType string `json:"msg_type"`
	util.LogConfigParams
}

