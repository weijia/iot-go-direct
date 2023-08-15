package msg

import "iot_go/pkg/bsp"

type BroadcastUpdateGlassColorRequest struct {
	Method string      `json:"method"`
	Params ColorParams `json:"params"`
}
type ColorParams struct {
	ColorParams string `json:"color"`
}

func (request BroadcastUpdateGlassColorRequest) handler() interface{} {
	reply := GatewayNodeIdReply{
		MsgType:       "gateway_reboot_reply",
		GatewayNodeID: bsp.BspConfigInstance.GatewayNodeID,
	}
	return reply
}
