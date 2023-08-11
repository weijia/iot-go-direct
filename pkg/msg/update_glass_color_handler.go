package msg

import (
	"iot_go/pkg/bsp"
)

type UpdateGlassColorRequest struct {
	Method string                   `json:"method"`
	Params []UpdateGlassColorParams `json:"params"`
}

type UpdateGlassColorParams struct {
	NodeID string `json:"node_id"`
	Color  string `json:"color"`
}

func (request UpdateGlassColorRequest) handle() interface{} {
	var reply UpdateGlassColorReply

	reply.GatewayNodeID = bsp.BspConfigInstance.GatewayNodeID
	reply.MsgType = "update_glass_color_reply"
	return reply
}
