package msg

type GetGlassStatusRequest struct {
	Method string            `json:"method"`
	Params GlassStatusParams `json:"params"`
}
type GlassStatusParams struct {
	NodeId string `json:"node_id"`
}
