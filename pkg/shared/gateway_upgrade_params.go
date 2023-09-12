package shared

type GatewayUpgradeParams struct {
	GatewayNodeId         string `json:"gateway_node_id"`
	TargetHardwareVersion string `json:"target_hardware_version"`
	TargetSoftwareVersion string `json:"target_software_version"`
	NodeType              int    `json:"node_type"`
	Crc8                  int    `json:"crc8"`
	IP                    string `json:"ip"`
	Port                  int    `json:"port"`
	User                  string `json:"user"`
	Pwd                   string `json:"pwd"`
	Path                  string `json:"path"`
}
