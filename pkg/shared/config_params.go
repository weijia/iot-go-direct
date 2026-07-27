package shared

// This is used to save to config file and receive server config from server
type BaseConfigParams struct {
	NodeList1      []string `json:"node_list_1"`
	NodeList2      []string `json:"node_list_2"`
	TouchNodeList1 []string `json:"touch_node_list_1"`
	TouchNodeList2 []string `json:"touch_node_list_2"`
	Custom         string   `json:"custom,omitempty"`
	Project        string   `json:"project,omitempty"`
	Heartbeat      int      `json:"heart_beat"`
	// Only accept multiples of 10, see usage of HeartbeatToServer
	HeartbeatToServer      int      `json:"heart_beat_to_server"`
	// SerialPollInterval 是网关串口层自动轮询颜色控制板状态的间隔（秒），
	// 与心跳间隔互相独立。控制板完成变色后由该轮询自动回填 completion_status，
	// 使心跳上报真实反映完成态，无需服务器手动发 get_glass_status。
	SerialPollInterval     int      `json:"serial_poll_interval"`
}

type ConfigParams struct {
	BaseConfigParams
	Module1 Module `json:"module1"`
	Module2 Module `json:"module2"`
}
