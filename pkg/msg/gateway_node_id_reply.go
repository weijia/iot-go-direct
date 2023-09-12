package msg

type GatewayNodeIdReply struct {
	MsgType       string `json:"msg_type"`
	GatewayNodeId string `json:"gateway_node_id"`
}
