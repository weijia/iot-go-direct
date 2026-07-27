package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
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
	var reply shared.GlassStatusUpdate
	reply.MsgType = "glass_status_update"
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId

	if bsp.IsInNodeList1(req.Params.NodeId) || bsp.IsInNodeList2(req.Params.NodeId) {
		if frame, ok := controller.Ctrl.QueryStatus(); ok {
			controller.Ctrl.UpdateNodeStatesFromSerialReply(frame)
		}
		nodeState := bsp.GetOrCreateNodeState(req.Params.NodeId)
		reply.Status = append(reply.Status, shared.GlassStatus{
			NodeId: req.Params.NodeId,
			Color:  GetColorStrFromSlice(nodeState.NodeReportedColor),
		})
	} else {
		reply.Status = append(reply.Status, shared.GlassStatus{
			NodeId: req.Params.NodeId,
			Color:  "1f2f3f4f5f6f7f8f",
		})
	}
	return reply
}
