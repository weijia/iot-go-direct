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
