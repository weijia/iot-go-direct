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
	case "node_list_request":
	 nodeListReply NodeListReply
		return nodeListReply.handle()
	case "node_info_request":
	 return nil
	case "gateway_reboot":
	 return BaseReply {"msg_type": "gateway_reboot_reply", 
	     "gateway_node_id": bsp.Bsp...
	 }
	case "update_glass_color_request":
	 
	}
	return nil
}
