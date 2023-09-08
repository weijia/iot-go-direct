package node

func GetHeartBeatMsg(gatewayId string, nodeId string) []byte {
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
