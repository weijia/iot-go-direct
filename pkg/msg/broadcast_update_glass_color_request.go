package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
)

type BroadcastUpdateGlassColorRequest struct {
	Method string      `json:"method"`
	Params ColorParams `json:"params"`
}
type ColorParams struct {
	ColorParams string `json:"color"`
}

func (request BroadcastUpdateGlassColorRequest) handle(mqttClient *util.Mqtt) interface{} {
	reply := GatewayNodeIdReply{
		MsgType:       "gateway_reboot_reply",
		GatewayNodeID: bsp.BspConfigInstance.GatewayNodeID,
	}
	return reply
}
