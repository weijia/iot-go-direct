package msg

type NodeFirmwareUpgradeReply struct {
	MsgType string `json:"msg_type"`
	// TargetHardwareVersion string `json:"target_hardware_version"`
	// TargetSoftwareVersion string `json:"target_software_version"`
	// NodeType              int    `json:"node_type"`
	// TargetRunArea         int    `json:"target_run_area"`
	// Crc8                  int    `json:"crc8"`
	// IP                    string `json:"ip"`
	// Port                  int    `json:"port"`
	NodeFirmwareUpgradeParams
	State string `json:"state"`
}
