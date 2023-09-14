package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
)

type GroupUpdateGlassColorRequest struct {
	Method string                        `json:"method"`
	Params []GroupUpdateGlassColorParams `json:"params"`
}
type GroupUpdateGlassColorParams struct {
	NodeList1 []string `json:"node_list1"`
	NodeList2 []string `json:"node_list2"`
	Color     string   `json:"color"`
}

func (request GroupUpdateGlassColorRequest) handle(mqttClient *util.Mqtt) interface{} {
	for _, group := range request.Params {
		for _, nodeId := range group.NodeList1 {
			client := bsp.GetLoraClientForNode(nodeId)
			if client != nil {
				SetSingleGlassColor(client, nodeId, group.Color)
			}
			
		}
		for _, nodeId := range group.NodeList2 {
			client := bsp.GetLoraClientForNode(nodeId)
			if client != nil {
				SetSingleGlassColor(client, nodeId, group.Color)
			}
			
		}
	}

	var reply GatewayNodeIdReply
	reply.MsgType = "group_update_glass_color_reply"
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	return reply
}
