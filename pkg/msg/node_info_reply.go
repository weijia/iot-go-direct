package msg

import (
	"iot_go/pkg/shared"
)

type NodeInfoReply struct {
	MsgType string `json:"msg_type"`
	// GatewayNodeID   string `json:"gateway_node_id"`
	// HardwareVersion string `json:"hardware_version"`
	// SoftwareVersion string `json:"software_version"`
	// Rssi            int    `json:"rssi"`
	// Ccid            string `json:"ccid"`
	// Custom          string `json:"custom"`
	// Project         string `json:"project"`
	// Module1 shared.Module `json:"module1"`
	// Module2 shared.Module `json:"module2"`
	shared.NodeInfoContent
	// MqttIP          string        `json:"mqtt_ip"`
	// MqttPort        int           `json:"mqtt_port"`
	// MqttUserName    string        `json:"mqtt_user_name"`
	// MqttPwd         string        `json:"mqtt_pwd"`
	shared.MqttParams
}
