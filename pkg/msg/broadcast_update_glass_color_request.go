package msg


type BroadcastUpdateGlassColorRequest struct {
	Method string `json:"method"`
	Params ColorParams `json:"params"`
}
type Params struct {
	ColorParams string `json:"color"`
}