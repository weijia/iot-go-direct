package msg

type NodeFirmwareUpgradeRequest struct {
	Method string                    `json:"method"`
	Params NodeFirmwareUpgradeParams `json:"params"`
}

type NodeFirmwareUpgradeParams struct {
	TargetHardwareVersion string `json:"target_hardware_version"`
	TargetSoftwareVersion string `json:"target_software_version"`
	NodeType              int    `json:"node_type"`
	TargetRunArea         int    `json:"target_run_area"`
	Crc8                  int    `json:"crc8"`
	IP                    string `json:"ip"`
	Port                  int    `json:"port"`
}

func (nodeFirmwareUpgradeRequest NodeFirmwareUpgradeRequest) handle() interface{} {
	var reply NodeFirmwareUpgradeReply
	reply.MsgType = "node_firmware_download_reply"
	reply.NodeFirmwareUpgradeParams = nodeFirmwareUpgradeRequest.Params
	return reply
}
