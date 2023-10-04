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
	reply.Status[0] = shared.GlassStatus{
		NodeId: nodeIdStr,
		Color:  "1f2f3f4f5f6f7f8f",
	}
	ch := make(chan node.NodeMsgReply)
	msgReq := node.NodeMsgReq{
		Data:    node.GetRetrieveColorMsg(nodeIdStr),
		ReplyCh: &ch,
	}
	for i := 0; i < util.GET_GLASS_STATUS_MAX_RETRY; i++ {
		loraModule.SendingCh <- msgReq
		n := node.GetReplyOrTimeout(ch)
		if !n.IsTimeout {
			reply.Status[0].Color = node.GetColorWithArea(n.Data)
			break
		}
	}
	mqttToServer <- reply
}
