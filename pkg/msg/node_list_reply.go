package msg

import (
	"iot_go/pkg/bsp"
)

type NodeListReply struct {
	MsgType        string   `json:"msg_type"`
	GatewayNodeID  string   `json:"gateway_node_id"`
	TotalCount     int      `json:"total_count"`
	NodeList1      []string `json:"node_list_1"`
	NodeList2      []string `json:"node_list_2"`
	TouchNodeList1 []string `json:"touch_node_list_1"`
	TouchNodeList2 []string `json:"touch_node_list_2"`
}

func (reply NodeListReply) handle() interface{} {
	reply.GatewayNodeID = bsp.BspConfigInstance.GatewayNodeID
	reply.NodeList1 = bsp.BspConfigInstance.NodeList1
	reply.NodeList2 = bsp.BspConfigInstance.NodeList2
	reply.TouchNodeList1 = bsp.BspConfigInstance.TouchNodeList1
	reply.TouchNodeList2 = bsp.BspConfigInstance.TouchNodeList2
	return reply
}
