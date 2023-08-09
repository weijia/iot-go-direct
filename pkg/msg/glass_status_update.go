package msg

type GlassStatusUpdate struct {
	MsgType       string        `json:"msg_type"`
	GatewayNodeID string        `json:"gateway_node_id"`
	Status        []GlassStatus `json:"status"`
}
type GlassStatus struct {
	NodeID string `json:"node_id"`
	Color  string `json:"color"`
}
