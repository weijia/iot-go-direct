package msg

type GetGlassStatusRequest struct {
	Method string            `json:"method"`
	Params GlassStatusParams `json:"params"`
}
type GlassStatusParams struct {
	NodeID string `json:"node_id"`
}
