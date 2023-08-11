package msg

import (
	"encoding/json"

	"iot_go/pkg/util"
)

type MsgHandler interface {
	handle() interface{}
}

type MsgFactory interface {
	getTargetMsg() any
}

func HandleMsg(body []byte) interface{} {
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
		return configRequest.handle()
	case "mqtt_config":
		var mqttConfigRequest MqttConfigRequest
		if err := json.Unmarshal(body, &mqttConfigRequest); err != nil {
			util.IotLogFatal(err)
		}
		return mqttConfigRequest.handle()
	}
	return nil
}
