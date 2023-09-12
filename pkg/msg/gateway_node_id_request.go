package msg

type GatewayNodeIdRequest struct {
	Method        string `json:"method"`
	GatewayNodeId string `json:"gateway_node_id"`
}
