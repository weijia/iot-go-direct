package msg

type GatewayNodeIdRequest struct {
	Method        string `json:"method"`
	GatewayNodeID string `json:"gateway_node_id"`
}
