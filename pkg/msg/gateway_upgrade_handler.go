package msg

import (
	"fmt"
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
	FIRMWARE_FOLDER    = "apps"
)

var IsRebootNeeded = false

// TODO: If we force the application to call another app even no other version exists, we can use this flag to upgrade without reboot
var IsQuitAppNeeded = false

func Download(params shared.GatewayUpgradeParams) error {

	targetPath := filepath.Join(util.GetOrCreateDownloadingFolder(), filepath.Base(params.SftpInfo.Path))
	err := util.Download(params.SftpInfo, targetPath, true)
	if err == nil {
		finalFilename := "main-" + params.TargetSoftwareVersion
		finalFullPath := filepath.Join(util.GetOrCreateAppFolder(), finalFilename)
		os.Rename(targetPath, finalFullPath)
		err := os.Chmod(finalFullPath, 0755)
		if err != nil {
			util.IotLogErrWithStr("Chmod error:", err)
		} else {
			IsRebootNeeded = true
		}
	}
	return err
}

func (req GatewayUpgradeRequest) handle() interface{} {
	var reply GatewayUpgradeReply
	reply.MsgType = "gateway_upgrade_reply"
	reply.GatewayUpgradeParams = req.Params
	reply.State = "OK"

	if req.Params.IsUpload != 0 {
		// Upload file
		err := util.Download(req.Params.SftpInfo, req.Params.SftpInfo.Path, false)
		if err != nil {
			reply.State = "error," + err.Error()
		}
		return reply
	}

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
		reply.State = fmt.Sprintf("error, target: %f is less than current version: %f", targetSwVer, localSwVer)
	} else {
		err := Download(req.Params)
		if err != nil {
			reply.State = "error," + err.Error()
		}
	}
	return reply
}
