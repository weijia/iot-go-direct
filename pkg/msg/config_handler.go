package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

type ConfigRequest struct {
	Method string              `json:"method"`
	Params shared.ConfigParams `json:"params"`
}

func (config ConfigRequest) handle(mqttClient *util.Mqtt) {
	node.HandleNodeInitReq(config.Params)
}

var IsConfigDone = false

func finalizeConfigReq(configParams shared.ConfigParams) GatewayNodeIdReply {
	bspInstance := bsp.GetBsp()
	bspInstance.SetModule1Params(configParams.Module1)
	bspInstance.SetModule2Params(configParams.Module2)

	bsp.BspConfigInstance.InitMsgContent.Module1 = configParams.Module1
	bsp.BspConfigInstance.InitMsgContent.Module2 = configParams.Module2
	bsp.BspConfigInstance.BaseConfigParams = configParams.BaseConfigParams
	bsp.BspConfigInstance.CommitChanges()
	var gatewayNodeIdReply GatewayNodeIdReply
	gatewayNodeIdReply.MsgType = "config_reply"
	gatewayNodeIdReply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	node.IsProcessingConfigReq = false
	return gatewayNodeIdReply
}
