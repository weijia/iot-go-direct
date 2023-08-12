package shared

type UpdateGlassColorReply struct {
	MsgType       string                   `json:"msg_type"`
	GatewayNodeID string                   `json:"gateway_node_id"`
	Status        []UpdateGlassColorParams `json:"status"`
}
