package msg

import (
	"encoding/json"

	"iot_go/pkg/util"
)

type MsgHandler interface {
	handle()
}

type MsgFactory interface {
	getTargetMsg() any
}

func HandleMsg(body []byte) {
	var baseRequest BaseRequest
	if err := json.Unmarshal(body, &baseRequest); err != nil {
		util.IotLogFatal(err)
	}
	switch baseRequest.Method {
	case "config":
		var configRequest ConfigRequest
		if err := json.Unmarshal(body, &configRequest); err != nil {
			util.IotLogFatal(err)
		}
		configRequest.handle()
	}
}
