package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

type ConfigRequest struct {
	Method string              `json:"method"`
	Params shared.ConfigParams `json:"params"`
}

func (config ConfigRequest) handle(mqttClient *util.Mqtt) interface{} {
	bspInstance := bsp.GetBsp()
	bspInstance.SetModule1Params(config.Params.Module1)
	bspInstance.SetModule2Params(config.Params.Module2)
	// fmt.Printf("%s", config.Method)
	bsp.BspConfigInstance.InitMsgContent.Module1 = config.Params.Module1
	bsp.BspConfigInstance.InitMsgContent.Module2 = config.Params.Module2
	bsp.BspConfigInstance.BaseConfigParams = config.Params.BaseConfigParams
	bsp.BspConfigInstance.CommitChanges()
	var gatewayNodeIdReply GatewayNodeIdReply
	gatewayNodeIdReply.MsgType = "config_reply"
	gatewayNodeIdReply.GatewayNodeID = bsp.BspConfigInstance.GatewayNodeID
	return gatewayNodeIdReply
}
