package msg

type Init struct {
	MsgType       string  `json:"msg_type"`
	GatewayNodeID string  `json:"gateway_node_id"`
	HardVersion   string  `json:"hard_version"`
	SoftVersion   string  `json:"soft_version"`
	Custom        string  `json:"custom"`
	Project       string  `json:"project"`
	NodeType      int     `json:"node_type"`
	Rssi          int     `json:"rssi"`
	Ccid          string  `json:"ccid"`
	HeartBeat     int     `json:"heart_beat"`
	Module0       Module0 `json:"module0"`
	Module1       Module1 `json:"module1"`
	Module2       Module2 `json:"module2"`
}
type Module0 struct {
	Freq   int `json:"freq"`
	Band   int `json:"band"`
	Factor int `json:"factor"`
}
type Module1 struct {
	Freq   int `json:"freq"`
	Band   int `json:"band"`
	Factor int `json:"factor"`
}
type Module2 struct {
	Freq   int `json:"freq"`
	Band   int `json:"band"`
	Factor int `json:"factor"`
}
