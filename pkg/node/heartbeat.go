package node

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
	"time"
)

func GetHeartBeatMsg(gatewayId []byte, nodeId []byte) []byte {
	var result []byte
	result = append(result, 0)            // package len
	result = append(result, 1)            // node type 1 gateway
	result = append(result, 1)            // cmd type 1 heartbeat
	result = append(result, gatewayId...) // gateway id
	result = append(result, nodeId...)    // node id
	result[0] = byte(len(result) + 1)     // need to count CRC byte in len
	result = append(result, getCRC8HighByTable(result))
	return result
}

func SendHeartbeat() {
	for {
		for _, value := range bsp.BspConfigInstance.BaseConfigParams.NodeList1 {
			bsp.GetModule1Client().Send(
				GetHeartBeatMsg(
					util.DecodeId(bsp.BspConfigInstance.GatewayNodeId), util.DecodeId(value)))
		}
		for _, value := range bsp.BspConfigInstance.BaseConfigParams.NodeList2 {
			bsp.GetModule2Client().Send(
				GetHeartBeatMsg(
					util.DecodeId(bsp.BspConfigInstance.GatewayNodeId), util.DecodeId(value)))
		}
		// Wait for specified time and send again
		time.Sleep(time.Duration(bsp.BspConfigInstance.BaseConfigParams.HeartBeat * 1000 * 1000 * 1000 * 60))
	}
}
