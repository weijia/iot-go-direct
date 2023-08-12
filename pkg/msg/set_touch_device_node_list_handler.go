package msg

import "fmt"

type SetTouchDeviceNodeList struct {
	Method string                `json:"method"`
	Params []TouchDeviceNodeList `json:"params"`
}
type TouchDeviceNodeList struct {
	NodeID  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}

func (request SetTouchDeviceNodeList) handle() interface{} {
	for _, value := range request.Params {
		fmt.Println(value.NodeID)
	}
	return nil
}
