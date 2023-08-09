package msg

type GatewayRebootReply struct {
	MsgType       string `json:"msg_type"`
	GatewayNodeID string `json:"gateway_node_id"`
}
