package msg

import (
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

type GatewayUpgradeRequest struct {
	Method string                      `json:"method"`
	Params shared.GatewayUpgradeParams `json:"params"`
}

func (gatewayUpgradeRequest GatewayUpgradeRequest) handle(mqttClient *util.Mqtt) interface{} {
	var reply GatewayUpgradeReply
	reply.MsgType = "node_firmware_download_reply"
	reply.GatewayUpgradeParams = gatewayUpgradeRequest.Params
	return reply
}
