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

func IsModuleParamChanged(config ConfigRequest) bool {
	if config.Params.Module1.Band != bsp.BspConfigInstance.Module1.Band ||
		config.Params.Module1.Freq != bsp.BspConfigInstance.Module1.Freq ||
		config.Params.Module1.Factor != bsp.BspConfigInstance.Module1.Factor ||
		config.Params.Module2.Band != bsp.BspConfigInstance.Module2.Band ||
		config.Params.Module2.Freq != bsp.BspConfigInstance.Module2.Freq ||
		config.Params.Module2.Factor != bsp.BspConfigInstance.Module2.Factor {
		return true
	}
	return false
}

func IsStrInSlice(str string, strSlice []string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}
	return false
}

// func AreTwoStrSliceEqual(first []string, second []string) bool {
// 	if len(first) == 0 && len(second) == 0 {
//         return true
//     }
//     if len(first) == 0 {
//         return false
//     }
//     if len(second) == 0 {
//         return false
//     }
//     for _, s := range second {
//         if IsStrInSlice(s, first) == false {
//             return false
//         }
//     }
//     for _, s := range first {
//         if IsStrInSlice(s, second) == false {
//             return false
//         }
//     }
//     return true
// }

// func IsNodeListChanged(config ConfigRequest) bool {
// 	if AreTwoStrSliceEqual(config.Params.NodeList1, bsp.BspConfigInstance.NodeList1) &&
// 			AreTwoStrSliceEqual(config.Params.NodeList2, bsp.BspConfigInstance.NodeList2) {
//         return false
//     }
//     return true
// }

var IsInitDone = false
var IsInitOngoing = false

func (config ConfigRequest) handle(mqttClient *util.Mqtt) interface{} {

	if IsInitDone {
		if IsModuleParamChanged(config) {
			// Need to send node init req and wait for resp in before freq change
			node.SendNodeInitReq(config.Params)
		} else {
			// Send node init req for newly added nodes

		}
	} else {
		if IsInitOngoing {
			//Ignore config if init config is already ongoing
			return nil
		} else {
			IsInitOngoing = true
			// Send heartbeat to nodes once
			node.SendHeartbeatOnce()
			// Send node init to nodes that does not respond
			go node.SendNodeHeartbeatInLoop()
			// TODO: Add node init reply 40 msg handler
		}
	}

}

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
