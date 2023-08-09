package msg

type SetTouchDeviceNodeList struct {
	Method string                `json:"method"`
	Params []TouchDeviceNodeList `json:"params"`
}
type TouchDeviceNodeList struct {
	NodeID  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}
