package msg

type NodeFirmwareDownloadRequest struct {
	Method string               `json:"method"`
	Params NodeFirmwareDownload `json:"params"`
}
type NodeFirmwareDownload struct {
	TargetHardwareVersion string `json:"target_hardware_version"`
	TargetSoftwareVersion string `json:"target_software_version"`
	NodeType              int    `json:"node_type"`
	TargetRunArea         int    `json:"target_run_area"`
	Crc8                  int    `json:"crc8"`
	IP                    string `json:"ip"`
	Port                  int    `json:"port"`
	User                  string `json:"user"`
	Pwd                   string `json:"pwd"`
	Path                  string `json:"path"`
}
