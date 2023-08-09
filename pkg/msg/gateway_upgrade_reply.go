package msg

type GatewayUpgradeReply struct {
	MsgType         string `json:"msg_type"`
	GatewayNodeID   string `json:"gateway_node_id"`
	HardwareVersion string `json:"hardware_version"`
	SoftwareVersion string `json:"software_version"`
	NodeType        int    `json:"node_type"`
	RunArea         int    `json:"run_area"`
	IP              string `json:"ip"`
	Port            int    `json:"port"`
	User            string `json:"user"`
	Pwd             string `json:"pwd"`
	Path            string `json:"path"`
	State           string `json:"state"`
}
