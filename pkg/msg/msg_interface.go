package msg

import (
	"encoding/json"

	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
)

type MsgHandler interface {
	handle() interface{}
}

type MsgFactory interface {
	getTargetMsg() any
}

func HandleMsg(mqttClient *util.Mqtt, body []byte) interface{} {
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
		var nodeListReply NodeListReply
		return nodeListReply.handle()
	case "node_info_request":
		return nil
	case "gateway_reboot":
		reply := GatewayNodeIdReply{
			MsgType:       "gateway_reboot_reply",
			GatewayNodeID: bsp.BspConfigInstance.GatewayNodeID,
		}
		return reply
	case "update_glass_color_request":
		var request UpdateGlassColorRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle(mqttClient)
	}
	return nil
}
