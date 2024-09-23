package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)

type BroadcastUpdateGlassColorRequest struct {
	Method string      `json:"method"`
	Params ColorParams `json:"params"`
}
type ColorParams struct {
	ColorParams string `json:"color"`
}

func (request BroadcastUpdateGlassColorRequest) handle() interface{} {
	if request.Params.ColorParams != "" {
		broadcastMsg := node.GetBroadcastUpdateGlassColorMsg(
			util.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
			request.Params.ColorParams)
		lora_module.Module1.SendWithoutReply(broadcastMsg)
		lora_module.Module2.SendWithoutReply(broadcastMsg)
		bsp.SetAllRequestingColor(request.Params.ColorParams)
	} else {
		util.IotLogErrorWithFormatStr("ColorParam is invalid: %v", request.Params.ColorParams)
	}
	reply := GatewayNodeIdReply{
		MsgType:       "broadcast_update_glass_color_reply",
		GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
	}
	return reply
}
