package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)

type GroupUpdateGlassColorRequest struct {
	Method string                      `json:"method"`
	Params GroupUpdateGlassColorParams `json:"params"`
}
type GroupUpdateGlassColorParams struct {
	NodeList1 []string `json:"node_list1"`
	NodeList2 []string `json:"node_list2"`
	Color     string   `json:"color"`
}

func (request GroupUpdateGlassColorRequest) Replay() {
	if GroupUpdateGlassColorRequestResendCnt > 0 {
		util.IotLogErrorWithFormatStr("resend group update glass color request, cnt: %d", GroupUpdateGlassColorRequestResendCnt)
		GroupUpdateGlassColorRequestResendCnt -= 1
		request.handle()
	}
}


func (request GroupUpdateGlassColorRequest) handle() interface{} {
	invalidNodeList := []string{}
	finalNodeList1 := [][]byte{}
	finalNodeList2 := [][]byte{}

	group := request.Params

	for _, nodeId := range group.NodeList1 {
		if bsp.IsInNodeList1(nodeId) {
			bsp.SetRequestingColor(nodeId, group.Color)
			finalNodeList1 = append(finalNodeList1, util.DecodeId(nodeId))
		} else {
			invalidNodeList = append(invalidNodeList, nodeId)
		}
	}
	if len(finalNodeList1) > 0 {
		groupMsg := node.GetGroupUpdateGlassColorMsg(
			finalNodeList1, group.Color)
		// bsp.GetModule1Client().Send(groupMsg)
		lora_module.Module1.SendWithoutReply(groupMsg)
	}

	for _, nodeId := range group.NodeList2 {
		if bsp.IsInNodeList2(nodeId) {
			bsp.SetRequestingColor(nodeId, group.Color)
			finalNodeList2 = append(finalNodeList2, util.DecodeId(nodeId))
		} else {
			invalidNodeList = append(invalidNodeList, nodeId)
		}
	}
	if len(finalNodeList2) > 0 {
		groupMsg := node.GetGroupUpdateGlassColorMsg(
			finalNodeList2, group.Color)
		// bsp.GetModule2Client().Send(groupMsg)
		lora_module.Module2.SendWithoutReply(groupMsg)
	}

	var reply GroupUpdateGlassColorReply
	reply.MsgType = "group_update_glass_color_reply"
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.InvalidNodes = invalidNodeList
	return reply
}
