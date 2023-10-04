package node

import (
	"iot_go/pkg/bsp"
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

// func waitUntilReplyOrTimeout(nodeIdStr string) {
// 	heartbeatTimer.Reset(util.HEARTBEAT_REPLY_TIMEOUT * time.Second)
// 	for {
// 		select {
// 		case <-heartbeatTimer.C:
// 			util.IotLogInfo(fmt.Sprintf("Heartbeat timeout for: %s\n", nodeIdStr))
// 			return
// 		case heartbeatReplyNodeIdStr := <-HeartbeatCh:
// 			util.IotLogInfo(fmt.Sprintf("Heartbeat reply for: %s\n", heartbeatReplyNodeIdStr))
// 			if nodeIdStr == heartbeatReplyNodeIdStr {
// 				return
// 			}
// 		}
// 	}
// }

// func sendHeartbeatForNodeList(client *lora_client.LoraClient, nodeList []string) {
// 	for _, nodeIdStr := range nodeList {
// 		util.IotLogInfo(fmt.Sprintf("Sending heartbeat for: %s", nodeIdStr))
// 		for i := 0; i < HeartbeatRetryCnt; i++ {
// 			client.Send(GetHeartBeatMsg(nodeIdStr))
// 			if util.IsReplyTimeout(nodeIdStr, HeartbeatCh, util.HEARTBEAT_REPLY_TIMEOUT) {
// 				util.IotLogInfo(fmt.Sprintf("Heartbeat reply timeout for: %s", nodeIdStr))
// 			} else {
// 				break
// 			}
// 		}
// 	}
// }

// var HeartbeatStartTime int64

// func SendHeartbeatOnce() {
// 	util.IotLogInfo("Sending a round of heartbeats")
// 	HeartbeatStartTime = time.Now().Unix()
// 	sendHeartbeatForNodeList(bsp.GetModule1Client(), bsp.BspConfigInstance.BaseConfigParams.NodeList1)
// 	sendHeartbeatForNodeList(bsp.GetModule2Client(), bsp.BspConfigInstance.BaseConfigParams.NodeList2)
// }

// func SendNodeHeartbeatInLoop() {
// 	HeartbeatRetryCnt = 3
// 	ticker1 := time.NewTicker(time.Duration(bsp.BspConfigInstance.Heartbeat) * time.Second)
// 	for {
// 		<-ticker1.C
// 		SendHeartbeatOnce()
// 	}
// }


type NodeMsgReply struct {
	Data []byte
	IsTimeout bool
}

type NodeMsgReq struct {
	Data []byte
	ReplyCh *chan NodeMsgReply
}

func GetReplyOrTimeout(ch chan NodeMsgReply) NodeMsgReply {
	eventTimer := time.NewTimer(time.Duration(20) * time.Second)
	reply := NodeMsgReply{
		Data: nil,
		IsTimeout: true,
	}
	select {
	case reply = <- ch:
		eventTimer.Stop()
	case <-eventTimer.C:
		util.IotLogErrorStr("Timeout for waiting for Module to reply")
	}
	return reply
}

func SendHeartbeatForNode(nodeIdStr string, sendingCh *chan NodeMsgReq) NodeMsgReply {
	ch := make(chan NodeMsgReply)
	msgReq := NodeMsgReq{
		Data: GetHeartBeatMsg(nodeIdStr),
		ReplyCh: &ch,
	}
	*sendingCh <- msgReq
	return GetReplyOrTimeout(ch)
}