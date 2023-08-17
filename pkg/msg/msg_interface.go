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

func GetMsgVar(method string) interface {} {
	var requestMap := map[string]interface{}{
		"node_firmware_upgrade_request": NodeFirmwareUpgradeRequest{},
		"config": ConfigRequest{},
	}
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
		var request NodeInfoRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle()
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
	case "group_update_glass_color_request":
		var request GroupUpdateGlassColorRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle(mqttClient)
	case "broadcast_update_glass_color_request":
		var request BroadcastUpdateGlassColorRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle()
	case "set_touch_device_node_list_request":
		var request SetTouchDeviceNodeList
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		request.handle()
		return nil
	case "gateway_upgrade_reply":
		var request GateWayUpgradeRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle()
	case "node_firmware_download_request":
		var request NodeFirmwareDownloadRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle()
	case "node_firmware_upgrade_request":
		var request NodeFirmwareUpgradeRequest
		if err := json.Unmarshal(body, &request); err != nil {
			util.IotLogFatal(err)
		}
		return request.handle()
	}
	return nil
}
