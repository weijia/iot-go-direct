package msg

type UpdateGlassColorReply struct {
	MsgType       string   `json:"msg_type"`
	GatewayNodeID string   `json:"gateway_node_id"`
	Status        []Status `json:"status"`
}
type Status struct {
	NodeID string `json:"node_id"`
	Color  string `json:"color"`
}
