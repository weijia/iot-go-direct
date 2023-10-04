package shared

type GlassStatusUpdate struct {
	MsgType       string        `json:"msg_type"`
	GatewayNodeId string        `json:"gateway_node_id"`
	Status        []GlassStatus `json:"status"`
}
type GlassStatus struct {
	NodeId string `json:"node_id"`
	Color  string `json:"color"`
}
