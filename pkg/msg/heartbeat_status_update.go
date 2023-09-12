package msg

type HeartbeatStatusUpdate struct {
	MsgType       string            `json:"msg_type"`
	GatewayNodeId string            `json:"gateway_node_id"`
	Status        []HeartbeatStatus `json:"status"`
}
type HeartbeatStatus struct {
	NodeId      string `json:"node_id"`
	Color       string `json:"color"`
	HardVersion string `json:"hard_version"`
	SoftVersion string `json:"soft_version"`
	RunArea     int    `json:"run_area"`
}
