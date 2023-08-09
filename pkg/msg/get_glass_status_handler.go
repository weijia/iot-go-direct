package msg

type GetGlassStatusRequest struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}
type Params struct {
	NodeID string `json:"node_id"`
}
