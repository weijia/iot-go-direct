package shared

type InitMsgContent struct {
	GatewayNodeID string `json:"gateway_node_id"`
	HardVersion   string `json:"hard_version"`
	SoftVersion   string `json:"soft_version"`
	Custom        string `json:"custom"`
	Project       string `json:"project"`
	NodeType      int    `json:"node_type"`
	Rssi          int    `json:"rssi"`
	Ccid          string `json:"ccid"`
	HeartBeat     int    `json:"heart_beat"`
	Module0       Module `json:"module0"`
	Module1       Module `json:"module1"`
	Module2       Module `json:"module2"`
}
