package lora_module

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

func (loraModule LoraModule) InitiateGetGlassStatusReq(
	nodeIdStr string, mqttToServer chan interface{}) {
	var reply shared.GlassStatusUpdate
	reply.MsgType = "glass_status_update"
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.Status = append(reply.Status, shared.GlassStatus{
		NodeId: nodeIdStr,
		Color:  "1f2f3f4f5f6f7f8f",
	})

	data := loraModule.SendNodeMsgWithRetryOrTimeout(node.GetRetrieveColorMsg(nodeIdStr), util.GET_GLASS_STATUS_MAX_RETRY)
	if data != nil {
		reply.Status[0].Color = node.GetColorWithArea(data)
	}
	mqttToServer <- reply
}
