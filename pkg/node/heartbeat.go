package node

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/util"
	"time"
)

func GetHeartBeatMsg(nodeIdStr string) []byte {
	var result []byte
	gatewayId := util.DecodeId(bsp.BspConfigInstance.GatewayNodeId)
	nodeId := util.DecodeId(nodeIdStr)
	result = append(result, 0)            // package len
	result = append(result, 1)            // node type 1 gateway
	result = append(result, 1)            // cmd type 1 heartbeat
	result = append(result, gatewayId...) // gateway id
	result = append(result, nodeId...)    // node id
	result[0] = byte(len(result) + 1)     // need to count CRC byte in len
	result = append(result, getCRC8HighByTable(result))
	return result
}

var HeartbeatRetryCnt = 1
var heartbeatTimer = time.NewTimer(util.HEARTBEAT_REPLY_TIMEOUT * time.Second)
var HeartbeatCh = make(chan string)

func waitUntilReplyOrTimeout(nodeIdStr string) {
	heartbeatTimer.Reset(util.HEARTBEAT_REPLY_TIMEOUT * time.Second)
	for {
		select {
		case <-heartbeatTimer.C:
			util.IotLogInfo(fmt.Sprintf("Heartbeat timeout for: %s\n", nodeIdStr))
			return
		case heartbeatReplyNodeIdStr := <-HeartbeatCh:
			util.IotLogInfo(fmt.Sprintf("Heartbeat reply for: %s\n", heartbeatReplyNodeIdStr))
			if nodeIdStr == heartbeatReplyNodeIdStr {
				return
			}
		}
	}
}

func sendHeartbeatForNodeList(client *lora_client.LoraClient, nodeList []string) {
	for _, nodeIdStr := range nodeList {
		util.IotLogInfo(fmt.Sprintf("Sending heartbeat for: %s\n", nodeIdStr))
		for i := 0; i < HeartbeatRetryCnt; i++ {
			client.Send(GetHeartBeatMsg(nodeIdStr))
			waitUntilReplyOrTimeout(nodeIdStr)
		}
	}
}

func SendHeartbeatOnce() {
	util.IotLogInfo("Sending a round of heartbeats\n")
	sendHeartbeatForNodeList(bsp.GetModule1Client(), bsp.BspConfigInstance.BaseConfigParams.NodeList1)
	sendHeartbeatForNodeList(bsp.GetModule2Client(), bsp.BspConfigInstance.BaseConfigParams.NodeList2)
}

func SendNodeHeartbeatInLoop() {
	HeartbeatRetryCnt = 3
	ticker1 := time.NewTicker(time.Duration(bsp.BspConfigInstance.HeartBeat) * time.Second)
	for {
		<-ticker1.C
		SendHeartbeatOnce()
	}
}
