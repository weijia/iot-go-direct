package node

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

// TODO: append prefix 0 to gatewayId and nodeId
func GetNodeInitMsg(gatewayId []byte, nodeId []byte, moduleParam shared.Module) []byte {
	var result []byte
	result = append(result, 0)            // package len
	result = append(result, 1)            // node type 1 gateway
	result = append(result, 3)            // cmd type 1 node init
	result = append(result, gatewayId...) // gateway id
	result = append(result, nodeId...)    // nod

	result = append(result, 3) // param num

	result = append(result, 1) // param type 1 freq
	result = append(result, byte(moduleParam.Freq>>8))
	result = append(result, byte(moduleParam.Freq&0xff))

	result = append(result, 2) // param type 2 factor
	result = append(result, byte(moduleParam.Factor>>8))
	result = append(result, byte(moduleParam.Factor&0xff))

	result = append(result, 3) // param type 3 bw
	result = append(result, byte(moduleParam.Band>>8))
	result = append(result, byte(moduleParam.Band&0xff))

	result[0] = byte(len(result) + 1) // need to count CRC byte in len
	result = append(result, getCRC8HighByTable(result))
	return result
}

func SendNodeInit(client *lora_client.LoraClient, nodeId string, moduelParam shared.Module) {
	util.IotLogInfo(fmt.Sprintf("Sending node init msg to %s with param: %v", nodeId, moduelParam))
	client.Send(GetNodeInitMsg(
		util.DecodeId(bsp.BspConfigInstance.GatewayNodeId), util.DecodeId(nodeId), moduelParam))
}

// func SendNodeInitAfterStartup() {
// 	for _, value := range bsp.BspConfigInstance.BaseConfigParams.NodeList1 {
// 		SendNodeInit(bsp.GetModule0Client(), value, bsp.BspConfigInstance.InitMsgContent.Module1)
// 	}
// 	for _, value := range bsp.BspConfigInstance.BaseConfigParams.NodeList2 {
// 		SendNodeInit(bsp.GetModule0Client(), value, bsp.BspConfigInstance.InitMsgContent.Module1)
// 	}
// }

var InitReplyCh = make(chan string)

func SendNodeInitReq(configParam shared.ConfigParams) {
	// For node already in node list, send init msg in working freq, otherwise send it in public freq
	for _, nodeId := range configParam.NodeList1 {
		client := bsp.GetLoraClientForNode(nodeId)
		if client != nil {
			SendNodeInit(client, nodeId, configParam.Module1)
			util.IsReplyTimeout(nodeId, InitReplyCh, 5)
		} else {
			SendNodeInit(bsp.GetModule0Client(), nodeId, configParam.Module1)
			util.IsReplyTimeout(nodeId, InitReplyCh, 5)
		}
	}
	for _, nodeId := range configParam.NodeList2 {
		client := bsp.GetLoraClientForNode(nodeId)
		if client != nil {
			SendNodeInit(client, nodeId, configParam.Module2)
			util.IsReplyTimeout(nodeId, InitReplyCh, 5)
		} else {
			SendNodeInit(bsp.GetModule0Client(), nodeId, configParam.Module2)
			util.IsReplyTimeout(nodeId, InitReplyCh, 5)
		}
	}
}
