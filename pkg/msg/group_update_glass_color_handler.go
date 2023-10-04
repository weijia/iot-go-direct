package msg

import (
	"iot_go/pkg/bsp"
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

func (request GroupUpdateGlassColorRequest) handle() interface{} {
	invalidNodeList := []string{}
	finalNodeList1 := [][]byte{}
	finalNodeList2 := [][]byte{}

	group := request.Params

	for _, nodeId := range group.NodeList1 {
		if bsp.IsInNodeList1(nodeId) {
			finalNodeList1 = append(finalNodeList1, util.DecodeId(nodeId))
		} else {
			invalidNodeList = append(invalidNodeList, nodeId)
		}
	}
	if len(finalNodeList1) > 0 {
		groupMsg := node.GetGroupUpdateGlassColorMsg(
			util.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
			finalNodeList1, group.Color)
		bsp.GetModule1Client().Send(groupMsg)
	}

	for _, nodeId := range group.NodeList2 {
		if bsp.IsInNodeList2(nodeId) {
			finalNodeList2 = append(finalNodeList2, util.DecodeId(nodeId))
		} else {
			invalidNodeList = append(invalidNodeList, nodeId)
		}
	}
	if len(finalNodeList2) > 0 {
		groupMsg := node.GetGroupUpdateGlassColorMsg(
			util.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
			finalNodeList2, group.Color)
		bsp.GetModule2Client().Send(groupMsg)
	}

	var reply GroupUpdateGlassColorReply
	reply.MsgType = "group_update_glass_color_reply"
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.InvalidNodes = invalidNodeList
	return reply
}
