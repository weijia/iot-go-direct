package msg

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"

	"github.com/xeipuuv/gojsonschema"
)

type MsgHandler interface {
	handle(*util.Mqtt) interface{}
}

//go:embed json_schema/update_color_schema.json
var update_color_schema string

//go:embed json_schema/init_msg_schema.json
var init_msg_schema string

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

func ValidateMsg(schemaStr string, data string) bool {
	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	documentLoader := gojsonschema.NewStringLoader(data)

	result, resultErr := gojsonschema.Validate(schemaLoader, documentLoader)
	if resultErr != nil {
		util.IotLogError(resultErr)
		return false
	} else {
		if result.Valid() {
			return true
		} else {
			for _, desc := range result.Errors() {
				util.IotLogErrorStr(fmt.Sprintf("- %s\n", desc))
			}
			return false
		}
	}
}

func HandleMqttMsg(mqttClient *util.Mqtt, body []byte) interface{} {
	var baseRequest BaseRequest
	if err := json.Unmarshal(body, &baseRequest); err != nil {
		util.IotLogError(err)
	}

	// targetMsg := GetMsgVar(baseRequest.Method).(MsgHandler)

	// if err := json.Unmarshal(body, &targetMsg); err != nil {
	// 	util.IotLogError(err)
	// }

	// if targetMsg != nil {
	// 	return targetMsg.(MsgHandler).handle(mqttClient)
	// }

	switch baseRequest.Method {

	case "broadcast_update_glass_color_request":
		var broadcastUpdateGlassColorRequest BroadcastUpdateGlassColorRequest
		if err := json.Unmarshal(body, &broadcastUpdateGlassColorRequest); err != nil {
			util.IotLogError(err)
		}
		return broadcastUpdateGlassColorRequest.handle(mqttClient)
	case "config":

		if !ValidateMsg(init_msg_schema, string(body)) {
			return nil
		}

		var configRequest ConfigRequest
		if err := json.Unmarshal(body, &configRequest); err != nil {
			util.IotLogError(err)
			return nil
		}
		configRequest.handle(mqttClient)
		return nil
	case "gateway_upgrade_reply":
		return nil

	case "group_update_glass_color_request":
		var req UpdateGlassColorRequest
		if err := json.Unmarshal(body, &req); err != nil {
			util.IotLogError(err)
		}
		req.handle()
		return nil

	case "node_info_request":
		var nodeInfoRequest NodeInfoRequest
		if err := json.Unmarshal(body, &nodeInfoRequest); err != nil {
			util.IotLogError(err)
		}
		return nodeInfoRequest.handle(mqttClient)
	case "node_list_request":
		var nodeListReply NodeListReply
		return nodeListReply.handle()
	case "gateway_reboot":
		reply := GatewayNodeIdReply{
			MsgType:       "gateway_reboot_reply",
			GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
		}
		return reply
	case "update_glass_color_request":

		if ValidateMsg(update_color_schema, string(body)) {
			return nil
		}

		var req UpdateGlassColorRequest

		if err := json.Unmarshal(body, &req); err != nil {
			util.IotLogError(err)
		}
		req.handle()
	}
	return nil
}
