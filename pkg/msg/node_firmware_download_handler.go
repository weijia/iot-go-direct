package msg


import (
	"iot_go/pkg/shared"
)

type NodeFirmwareDownloadRequest struct {
	Method string               `json:"method"`
	Params NodeFirmwareDownloadParam `json:"params"`
}

type NodeFirmwareDownloadParam struct {
	TargetHardwareVersion string `json:"target_hardware_version"`
	TargetSoftwareVersion string `json:"target_software_version"`
	NodeType              int    `json:"node_type"`
	TargetRunArea         int    `json:"target_run_area"`
	Crc8                  int    `json:"crc8"`
	shared.SftpInfo
	// IP                    string `json:"ip"`
	// Port                  int    `json:"port"`
	// User                  string `json:"user"`
	// Pwd                   string `json:"pwd"`
	// Path                  string `json:"path"`
}


func (nodeFirmwareDownloadRequest NodeFirmwareDownloadRequest) handle() interface {}{
	var reply NodeFirmwareDownloadReply
	reply.MsgType = "node_firmware_download_reply"
	reply.NodeFirmwareDownload= nodeFirmwareDownloadRequest.Params
	return reply
}