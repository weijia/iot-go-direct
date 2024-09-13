package shared

// Init message data from gateway to server
type NodeInfoContent struct {
	GatewayNodeId string `json:"gateway_node_id"`
	HardVersion   string `json:"hard_version"`
	SoftVersion   string `json:"soft_version"`
	Custom        string `json:"custom,omitempty"`
	Project       string `json:"project,omitempty"`
	NodeType      int    `json:"node_type"`
	Rssi          int    `json:"rssi"`
	Ccid          string `json:"ccid"`
	Heartbeat     int    `json:"heart_beat"`
	HeartbeatToServer     int    `json:"heart_beat_to_server"`
	Module1       Module `json:"module1"`
	Module2       Module `json:"module2"`
}

// Init message data from gateway to server
type InitMsgContent struct {
	NodeInfoContent
	Module0 Module `json:"module0"`
}
