package shared

// This is used to save to config file and receive server config from server
type BaseConfigParams struct {
	NodeList1      []string `json:"node_list_1"`
	NodeList2      []string `json:"node_list_2"`
	TouchNodeList1 []string `json:"touch_node_list_1"`
	TouchNodeList2 []string `json:"touch_node_list_2"`
	Custom         string   `json:"custom"`
	Project        string   `json:"project"`
	Heartbeat      int      `json:"heart_beat"`
	// Only accept multiples of 10, see usage of HeartbeatToServer
	HeartbeatToServer      int      `json:"heart_beat_to_server"`
}

type ConfigParams struct {
	BaseConfigParams
	Module1 Module `json:"module1"`
	Module2 Module `json:"module2"`
}
