package msg

import (
	"encoding/json"

	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
)

type MsgHandler interface {
	handle(*util.Mqtt) interface{}
}

func GetMsgVar(method string) interface{} {
	requestMap := map[string]interface{}{
		"broadcast_update_glass_color_request": BroadcastUpdateGlassColorRequest{},
		"config":                               ConfigRequest{},
		"gateway_upgrade_reply":                GatewayUpgradeRequest{},
		"group_update_glass_color_request":     GroupUpdateGlassColorRequest{},
		"mqtt_config":                          MqttConfigRequest{},
		"node_firmware_download_request":       NodeFirmwareDownloadRequest{},
		"node_firmware_upgrade_request":        NodeFirmwareUpgradeRequest{},
		"node_info_request":                    NodeInfoRequest{},
		"set_touch_device_node_list_request":   SetTouchDeviceNodeList{},
		"update_glass_color_request":           UpdateGlassColorRequest{},
	}
	return requestMap[method]
}

func HandleMsg(mqttClient *util.Mqtt, body []byte) interface{} {
	var baseRequest BaseRequest
	if err := json.Unmarshal(body, &baseRequest); err != nil {
		util.IotLogFatal(err)
	}

	targetMsg := GetMsgVar(baseRequest.Method)

	if targetMsg != nil {
		return targetMsg.(MsgHandler).handle(mqttClient)
	}

	switch baseRequest.Method {

	case "node_list_request":
		var nodeListReply NodeListReply
		return nodeListReply.handle()
	case "gateway_reboot":
		reply := GatewayNodeIdReply{
			MsgType:       "gateway_reboot_reply",
			GatewayNodeID: bsp.BspConfigInstance.GatewayNodeID,
		}
		return reply
	}
	return nil
}
