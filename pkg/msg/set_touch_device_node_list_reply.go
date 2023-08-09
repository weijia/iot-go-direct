package msg

type SetTouchDeviceNodeListReply struct {
	MsgType       string `json:"msg_type"`
	GatewayNodeID string `json:"gateway_node_id"`
	State         string `json:"state"`
}
