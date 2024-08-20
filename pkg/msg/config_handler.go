package msg

import (
	"context"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"sync"
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

func InitAccordingToConfig(ctx context.Context) {
	IsInitOngoing = true

	var wg sync.WaitGroup
	wg.Add(2)
	// Always use copy of the node list to prevent
	// accessing node list from different go routine
	copiedNodeList1 := make([]string, len(bsp.BspConfigInstance.NodeList1)) // 预先分配足够的空间  
    copy(copiedNodeList1, bsp.BspConfigInstance.NodeList1) // 复制元素  
	copiedNodeList2 := make([]string, len(bsp.BspConfigInstance.NodeList2)) // 预先分配足够的空间  
    copy(copiedNodeList2, bsp.BspConfigInstance.NodeList2) // 复制元素  

	// Send heartbeat to Module1
	go lora_module.Module1.SendHeartbeatForList(ctx, copiedNodeList1, &wg)
	// Send heartbeat to Module2
	go lora_module.Module2.SendHeartbeatForList(ctx, copiedNodeList2, &wg)
	// Not responding node will receive node init in public freq in above steps
	// Wait for above 2 activities to complete
	wg.Wait()

	IsInitOngoing = false
	IsInitDone = true
	go lora_module.Module1.SendHeartbeatForListInLoop(ctx, copiedNodeList1)
	go lora_module.Module2.SendHeartbeatForListInLoop(ctx, copiedNodeList2)
}

func (config ConfigRequest) handle(ctx context.Context) interface{} {
	if bsp.PeriodNumberForReportingToServer != config.Params.HeartbeatToServer/10 {
		bsp.PeriodNumberForReportingToServer = config.Params.HeartbeatToServer/10
		util.IotLog("Updated reporting period: %d seconds", bsp.PeriodNumberForReportingToServer*10)
	}
	
	if IsInitDone {
		go HandleConfigAfterInitWillBlock(config.Params)
		// TODO: if freq is not changed, we do not need to re init lora
		// if IsModuleParamChanged(config) {
		// Need to send node init req and wait for resp in before freq change
		return getConfigReply()
		// } else {
		// Send node init req for newly added nodes

		// }
	} else {
		// Init is not done
		if IsInitOngoing {
			//Ignore config if init config is already ongoing
			util.IotLogErrorStr("Receive init msg from server when init started already")
			return nil
		} else {
			saveConfig(config.Params)
			go InitAccordingToConfig(ctx)
			return getConfigReply()
		}
	}
}

func getConfigReply() interface{} {
	var gatewayNodeIdReply GatewayNodeIdReply
	gatewayNodeIdReply.MsgType = "config_reply"
	gatewayNodeIdReply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	// node.IsProcessingConfigReq = false
	return gatewayNodeIdReply
}

func saveConfig(configParams shared.ConfigParams) {
	// Update InitMsgContent so when gateway starts, latest config will be sent to server
	bsp.BspConfigInstance.InitMsgContent.Module1 = configParams.Module1
	bsp.BspConfigInstance.InitMsgContent.Module2 = configParams.Module2

	// We will only accept number multiple 10 as HeartbeatToServer, see HeartbeatToServer handling
	bsp.BspConfigInstance.InitMsgContent.HeartbeatToServer = 10*configParams.HeartbeatToServer/10

	bsp.BspConfigInstance.BaseConfigParams = configParams.BaseConfigParams
	if bsp.BspConfigInstance.BaseConfigParams.Heartbeat <=0 {
		bsp.BspConfigInstance.BaseConfigParams.Heartbeat = 20
		bsp.BspConfigInstance.InitMsgContent.Heartbeat = bsp.BspConfigInstance.Heartbeat
	}
	// TODO: update work around for list 2 empty
	if len(bsp.BspConfigInstance.BaseConfigParams.NodeList2) > 0 && bsp.BspConfigInstance.BaseConfigParams.NodeList2[0] == "        " {
		bsp.BspConfigInstance.BaseConfigParams.NodeList2 = 
			bsp.BspConfigInstance.BaseConfigParams.NodeList2[:0]
	}
	bsp.BspConfigInstance.CommitChanges()
}

func HandleConfigAfterInitWillBlock(configParams shared.ConfigParams) {
	needSendInitForModule0 := make(map[string]shared.Module)
	needSendInitForModule1 := make(map[string]shared.Module)
	needSendInitForModule2 := make(map[string]shared.Module)

	// config node list1 may contain node that is currently for module 2 or not joined (send init in module0)
	for _, nodeIdStr := range configParams.NodeList1 {
		if bsp.IsInNodeList1(nodeIdStr) {
			needSendInitForModule1[nodeIdStr] = configParams.Module1
		} else if bsp.IsInNodeList2(nodeIdStr) {
			needSendInitForModule2[nodeIdStr] = configParams.Module1
		} else {
			needSendInitForModule0[nodeIdStr] = configParams.Module1
		}
	}
	for _, nodeIdStr := range configParams.NodeList2 {
		if bsp.IsInNodeList1(nodeIdStr) {
			needSendInitForModule1[nodeIdStr] = configParams.Module2
		} else if bsp.IsInNodeList2(nodeIdStr) {
			needSendInitForModule2[nodeIdStr] = configParams.Module2
		} else {
			needSendInitForModule0[nodeIdStr] = configParams.Module2
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go lora_module.Module0.SendNodeInitForListWillBlock(needSendInitForModule0, &wg)
	go lora_module.Module1.SendNodeInitForListWillBlock(needSendInitForModule1, &wg)
	go lora_module.Module2.SendNodeInitForListWillBlock(needSendInitForModule2, &wg)
	wg.Wait()

	saveConfig(configParams)

	lora_module.Module1.UpdateNodeList(bsp.BspConfigInstance.BaseConfigParams.NodeList1)
	lora_module.Module2.UpdateNodeList(bsp.BspConfigInstance.BaseConfigParams.NodeList2)

	bsp.GetBsp().SetModule1HWParams(configParams.Module1)
	bsp.GetBsp().SetModule2HWParams(configParams.Module2)
}
