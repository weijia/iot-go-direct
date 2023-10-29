package node

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
)

func GetHeartbeatMsg(nodeIdStr string) []byte {
	var result []byte
	gatewayId := util.DecodeId(bsp.BspConfigInstance.GatewayNodeId)
	nodeId := util.DecodeId(nodeIdStr)
	result = append(result, 0)            // package len
	result = append(result, 1)            // node type 1 gateway
	result = append(result, HEARTBEAT_REQ)            // cmd type 1 heartbeat
	result = append(result, gatewayId...) // gateway id
	result = append(result, nodeId...)    // node id
	result[0] = byte(len(result) + 1)     // need to count CRC byte in len
	result = append(result, getCRC8HighByTable(result))
	return result
}
