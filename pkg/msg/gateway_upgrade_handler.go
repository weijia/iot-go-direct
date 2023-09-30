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

func Download(params shared.SftpInfo) error {
	os.MkdirAll("firmware", 771)
	targetPath := "firmware/" + filepath.Base(params.Path)
	err := util.Download(params, targetPath)
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
		err := Download(req.Params.SftpInfo)
		if err != nil {
			reply.State = "error," + err.Error()
		} else {

		}
	}
	return reply
}
