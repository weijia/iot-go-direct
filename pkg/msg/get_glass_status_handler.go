package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/shared"
)

type GetGlassStatusRequest struct {
	Method string            `json:"method"`
	Params GlassStatusParams `json:"params"`
}
type GlassStatusParams struct {
	NodeId string `json:"node_id"`
}

func (req GetGlassStatusRequest) handle(mqttToServer chan interface{}) interface{} {
	if bsp.IsInNodeList1(req.Params.NodeId) {
		go lora_module.Module1.InitiateGetGlassStatusReq(req.Params.NodeId, mqttToServer)
	} else if bsp.IsInNodeList2(req.Params.NodeId) {
		go lora_module.Module1.InitiateGetGlassStatusReq(req.Params.NodeId, mqttToServer)
	} else {
		var reply shared.GlassStatusUpdate
		reply.MsgType = "glass_status_update"
		reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
		reply.Status[0] = shared.GlassStatus{
			NodeId: req.Params.NodeId,
			Color:  "1f2f3f4f5f6f7f8f",
		}
		return reply
	}
	return nil
}
