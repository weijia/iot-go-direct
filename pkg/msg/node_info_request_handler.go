package msg

type NodeInfoRequest struct {
	Method string        `json:"method"`
	Params GatewayParams `json:"params"`
}
type GatewayParams struct {
	GatewayNodeID string `json:"gateway_node_id"`
}
