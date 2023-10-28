package node

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

func GetNodeInitMsg(nodeIdStr string, moduleParam shared.Module) []byte {
	var result []byte
	result = append(result, 0)                                                     // package len
	result = append(result, 1)                                                     // node type 1 gateway
	result = append(result, CONFIG_NODE_REQ)                                       // cmd type
	result = append(result, util.DecodeId(bsp.BspConfigInstance.GatewayNodeId)...) // gateway id
	result = append(result, util.DecodeId(nodeIdStr)...)                           // node id

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
