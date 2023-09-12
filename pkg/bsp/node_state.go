package bsp

import "iot_go/pkg/shared"

type NodeState struct {
	NodeId            string `json:"node_id"`
	LastMsgTimestamp  int64  `json:"node_msg_timestamp"`
	NodeReportedColor []int  `json:"node_reported_color"`
	// NodeRequestingColor []int  `json:"node_requesting_color"`
	HwVer       int           `json:"hardware_version"`
	SwVer       int           `json:"software_version"`
	RunningArea int           `json:"running_area"`
	ModuleParam shared.Module `json:"module_param"`
}

func GetNodeState(nodeId string) *NodeState {
	for _, nodeState := range BspConfigInstance.NodeStates {
		if nodeState.NodeId == nodeId {
			return &nodeState
		}
	}
	return nil
}
