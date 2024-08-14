package bsp

import (
	"iot_go/pkg/shared"
)

type NodeState struct {
	NodeId            string `json:"node_id"`
	LastMsgTimestamp  int64  `json:"node_msg_timestamp"`
	NodeReportedColor [8]int `json:"node_reported_color"`
	// NodeRequestingColor []int  `json:"node_requesting_color"`
	HwVer       int           `json:"hardware_version"`
	SwVer       int           `json:"software_version"`
	RunningArea int           `json:"running_area"`
	ModuleParam shared.Module `json:"module_param"`
	RSSI float64 `json:"rssi"`
	SNR float64 `json:"snr"`
	IsOffline bool `json:"is_offline"`
	CompletionStatus int `json:"completion_status"`
}

func GetOrCreateNodeState(nodeId string) *NodeState {
	for index := range BspConfigInstance.NodeStates {
		if BspConfigInstance.NodeStates[index].NodeId == nodeId {
			return &BspConfigInstance.NodeStates[index]
		}
	}
	newNode := NodeState{
		NodeId: nodeId,
	}
	BspConfigInstance.NodeStates = append(BspConfigInstance.NodeStates, newNode)
	return &BspConfigInstance.NodeStates[len(BspConfigInstance.NodeStates)-1]
}
