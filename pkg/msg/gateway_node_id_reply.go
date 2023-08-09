package msg

type GatewayNodeIdReply struct {
	MsgType       string `json:"msg_type"`
	GatewayNodeID string `json:"gateway_node_id"`
}
