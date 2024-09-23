package msg

import (
	"fmt"
)

type HeartbeatStatusUpdate struct {
	MsgType       string            `json:"msg_type"`
	GatewayNodeId string            `json:"gateway_node_id"`
	Status        []HeartbeatStatus `json:"status"`
}
type HeartbeatStatus struct {
	NodeId      string `json:"node_id"`
	Color       string `json:"color"`
	HardVersion string `json:"hard_version"`
	SoftVersion string `json:"soft_version"`
	RunArea     int    `json:"run_area"`
	RSSI float64 `json:"rssi"`
	SNR float64 `json:"snr"`
	IsOffline bool `json:"is_offline"`
	CompletionStatus int `json:"completion_status"`
	NodeRequestingColor string        `json:"node_requesting_color"`
}

func GetColorStrFromSlice(colorSlice [8]int) string {
	res := ""
	for index, color := range colorSlice {
		res += fmt.Sprintf("%x%x", index+1, color)
	}
	return res
}
