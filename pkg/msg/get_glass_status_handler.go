package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)

type GetGlassStatusRequest struct {
	Method string            `json:"method"`
	Params GlassStatusParams `json:"params"`
}
type GlassStatusParams struct {
	NodeId string `json:"node_id"`
}

func (req GetGlassStatusRequest) handle() interface{} {
	msg := node.GetRetrieveColorMsg(util.DecodeId(req.Params.NodeId))
	client := bsp.GetLoraClientForNode(req.Params.NodeId)
	if client != nil {
		client.Send(msg)
		return nil
	} else {
		var reply GlassStatusUpdate
		reply.MsgType = "glass_status_update"
		reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
		reply.Status[0] = GlassStatus{
			NodeId: req.Params.NodeId,
			Color:  "1f2f3f4f5f6f7f8f",
		}
		return reply
	}
}
