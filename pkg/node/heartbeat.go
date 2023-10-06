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

type NodeMsgReply struct {
	Data      []byte
	IsTimeout bool
}

type NodeMsgReq struct {
	Data    []byte
	ReplyCh *chan NodeMsgReply
}

func GetReplyOrTimeout(ch *chan NodeMsgReply) NodeMsgReply {
	eventTimer := time.NewTimer(time.Duration(util.LEVEL2_NODE_MSG_REPLY_TIMEOUT_SECONDS) * time.Second)
	reply := NodeMsgReply{
		Data:      nil,
		IsTimeout: true,
	}
	select {
	case reply = <-*ch:
		eventTimer.Stop()
		// util.IotLog("GetReplyOrTimeout received reply from level1, returning: %v", reply)
	case <-eventTimer.C:
		util.IotLogErrorStr("Timeout for waiting for Module to reply")
	}
	return reply
}

func SendHeartbeatForNode(nodeIdStr string, sendingCh *chan NodeMsgReq) NodeMsgReply {
	ch := make(chan NodeMsgReply)
	msgReq := NodeMsgReq{
		Data:    GetHeartBeatMsg(nodeIdStr),
		ReplyCh: &ch,
	}
	*sendingCh <- msgReq
	// util.IotLog("After request to send heartbeat")
	return GetReplyOrTimeout(&ch)
}
