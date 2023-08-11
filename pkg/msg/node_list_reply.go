package msg

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
    reply.GatewayNodeID = bsp.BspConfigInstance.InitConfig.GatewayNodeID
    reply.NodeList1 = bsp.BspConfigInstance.InitConfig.NodeList1
    reply.NodeList2 = bsp.BspConfigInstance.InitConfig.NodeList2
    reply.TouchNodeList1 = bsp.BspConfigInstance.InitConfig.TouchNodeList1
    reply.TouchNodeList2 = bsp.BspConfigInstance.InitConfig.TouchNodeList2
    return reply
}