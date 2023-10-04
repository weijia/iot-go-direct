package msg

import (
	"iot_go/pkg/bsp"
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
	broadcastMsg := node.GetBroadcastUpdateGlassColorMsg(
		util.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		request.Params.ColorParams)
	bsp.GetModule1Client().Send(broadcastMsg)
	bsp.GetModule2Client().Send(broadcastMsg)
	reply := GatewayNodeIdReply{
		MsgType:       "gateway_reboot_reply",
		GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
	}
	return reply
}
