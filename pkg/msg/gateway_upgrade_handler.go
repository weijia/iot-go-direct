package msg

import (
	"os"
	"path/filepath"

	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"strconv"
)

type GatewayUpgradeRequest struct {
	Method string                      `json:"method"`
	Params shared.GatewayUpgradeParams `json:"params"`
}

const (
	DOWNLOADING_FOLDER = "downloading"
	FIRMWARE_FOLDER = "apps"
)

var isRebootNeeded = false

func Download(params shared.GatewayUpgradeParams) error {
	os.MkdirAll(DOWNLOADING_FOLDER, 0771)
	targetPath := filepath.Join(DOWNLOADING_FOLDER, filepath.Base(params.SftpInfo.Path))
	err := util.Download(params.SftpInfo, targetPath)
	if err == nil {
		finalFilename := "main-" + params.TargetSoftwareVersion
		os.Rename(targetPath, filepath.Join(FIRMWARE_FOLDER, finalFilename))
		os.Chmod(finalFilename, 0755)
		isRebootNeeded = true
	}
	return err
}

func (req GatewayUpgradeRequest) handle() interface{} {
	var reply GatewayUpgradeReply
	reply.MsgType = "gateway_upgrade_reply"
	reply.GatewayUpgradeParams = req.Params
	reply.State = "OK"
	targetSwVer, err := strconv.ParseFloat(req.Params.TargetSoftwareVersion, 32)
	if err != nil {
		reply.State = "error, " + err.Error()
		return reply
	}
	localSwVer, err := strconv.ParseFloat(bsp.BspConfigInstance.SoftVersion, 32)
	if err != nil {
		reply.State = "error, " + err.Error()
		return reply
	}
	if req.Params.GatewayNodeId !=
			bsp.BspConfigInstance.GatewayNodeId ||
			req.Params.TargetHardwareVersion !=
				bsp.BspConfigInstance.HardVersion ||
			targetSwVer <= localSwVer {
		reply.State = "error"
	} else {
		err := Download(req.Params)
		if err != nil {
			reply.State = "error," + err.Error()
		}
	}
	return reply
}
