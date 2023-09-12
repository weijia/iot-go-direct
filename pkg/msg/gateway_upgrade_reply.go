package msg

import (
	"iot_go/pkg/shared"
)

type GatewayUpgradeReply struct {
	MsgType string `json:"msg_type"`
	// GatewayNodeId   string `json:"gateway_node_id"`
	// HardwareVersion string `json:"hardware_version"`
	// SoftwareVersion string `json:"software_version"`
	// NodeType        int    `json:"node_type"`
	// IP              string `json:"ip"`
	// Port            int    `json:"port"`
	// User            string `json:"user"`
	// Pwd             string `json:"pwd"`
	// Path            string `json:"path"`
	shared.GatewayUpgradeParams
	RunArea int    `json:"run_area"`
	State   string `json:"state"`
}
