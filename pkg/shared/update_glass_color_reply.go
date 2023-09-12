package shared

type UpdateGlassColorReply struct {
	MsgType       string                   `json:"msg_type"`
	GatewayNodeId string                   `json:"gateway_node_id"`
	Status        []UpdateGlassColorParams `json:"status"`
}
