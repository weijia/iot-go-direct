package node

import "iot_go/pkg/shared"

func GetNodeInitMsg(gatewayId string, nodeId string, moduleParam shared.Module) []byte {
	var result []byte
	result = append(result, 0)            // package len
	result = append(result, 1)            // node type 1 gateway
	result = append(result, 3)            // cmd type 1 node init
	result = append(result, gatewayId...) // gateway id
	result = append(result, nodeId...)    // nod

	result = append(result, 3) // param num

	result = append(result, 1)                      // param type 1 freq
	result = append(result, byte(moduleParam.Freq)) // param type 1 freq

	result = append(result, 1)                        // param type 2 factor
	result = append(result, byte(moduleParam.Factor)) // param type 1 freq

	result = append(result, 1)                      // param type 3 bw
	result = append(result, byte(moduleParam.Band)) // param type 1 freq

	result[0] = byte(len(result) + 1) // need to count CRC byte in len
	result = append(result, getCRC8HighByTable(result))
	return result
}
