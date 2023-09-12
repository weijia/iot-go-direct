package msg

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"time"

	"github.com/xeipuuv/gojsonschema"
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

func ValidateMsg(schemaStr string, data string) bool {
	loader1 := gojsonschema.NewStringLoader(schemaStr)
	schema, schemaErr := gojsonschema.NewSchema(loader1)
	if schemaErr != nil {
		util.IotLogError(schemaErr)
		return false
	} else {
		result, resultErr := schema.Validate(gojsonschema.NewGoLoader(data))
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
}

type ColorUpdateState struct {
	Reply        shared.UpdateGlassColorReply
	RepliedNodes []string
	Timer        *time.Timer
}

var statesForPendingColorUpdate []ColorUpdateState

var colorUpdateRequestTimeoutCh = make(chan *shared.UpdateGlassColorReply)

var configReqCh = make(chan ConfigRequest)

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

		//go:embed json_schema/init_msg_schema.json
		var s string

		if ValidateMsg(s, string(body)) {
			return nil
		}

		var configRequest ConfigRequest
		if err := json.Unmarshal(body, &configRequest); err != nil {
			util.IotLogError(err)
			return nil
		}
		return configRequest.handle(mqttClient)
	case "gateway_upgrade_reply":
		return nil

	case "group_update_glass_color_request":
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
		//go:embed json_schema/update_color_schema.json
		var s string

		if ValidateMsg(s, string(body)) {
			return nil
		}

		var req UpdateGlassColorRequest

		if err := json.Unmarshal(body, &req); err != nil {
			util.IotLogError(err)
		}
		reply := req.handle(mqttClient).(shared.UpdateGlassColorReply)

		state := ColorUpdateState{
			Reply: reply,
			Timer: time.NewTimer(120 * time.Second),
		}

		statesForPendingColorUpdate = append(statesForPendingColorUpdate, state)
		// Send reply pointer to requestTimeoutCh after 120 seconds
		// TODO: Maybe we need to retry before actual 120 second timeout
		go func() {
			select {
			case <-state.Timer.C:
				colorUpdateRequestTimeoutCh <- &state.Reply
			}
		}()
		return nil
	}
	return nil
}
