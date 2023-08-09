package msg

type UpdateGlassColorRequest struct {
	Method string                   `json:"method"`
	Params []UpdateGlassColorParams `json:"params"`
}
type UpdateGlassColorParams struct {
	NodeID string `json:"node_id"`
	Color  string `json:"color"`
}
