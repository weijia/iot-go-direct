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

// We can only handle 1 config request at a time
var OngoingConfigReq ConfigRequest
var isProcessingConfigReq bool = false

func (config ConfigRequest) handle(mqttClient *util.Mqtt) interface{} {
	if isProcessingConfigReq {
		return nil
	}
	OngoingConfigReq = config
	isProcessingConfigReq = true

	// For node already in node list, send init msg in working freq
	for _, value := range config.Params.NodeList1 {
		client := bsp.GetLoraClientForNode(value)
		if client != nil {
			bsp.SendNodeInit(client, value, config.Params.Module1)
		} else {
			bsp.SendNodeInit(bsp.GetModule0Client(), value, config.Params.Module1)
		}
	}
	for _, value := range config.Params.NodeList2 {
		client := bsp.GetLoraClientForNode(value)
		if client != nil {
			bsp.SendNodeInit(client, value, config.Params.Module2)
		} else {
			bsp.SendNodeInit(bsp.GetModule0Client(), value, config.Params.Module2)
		}
	}
}

func finalizeConfigReq(config ConfigRequest) GatewayNodeIdReply {
	bspInstance := bsp.GetBsp()
	bspInstance.SetModule1Params(config.Params.Module1)
	bspInstance.SetModule2Params(config.Params.Module2)

	bsp.BspConfigInstance.InitMsgContent.Module1 = config.Params.Module1
	bsp.BspConfigInstance.InitMsgContent.Module2 = config.Params.Module2
	bsp.BspConfigInstance.BaseConfigParams = config.Params.BaseConfigParams
	bsp.BspConfigInstance.CommitChanges()
	var gatewayNodeIdReply GatewayNodeIdReply
	gatewayNodeIdReply.MsgType = "config_reply"
	gatewayNodeIdReply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	return gatewayNodeIdReply
}
