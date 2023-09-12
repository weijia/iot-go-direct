package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"time"
)

const (
	HEARTBEAT_REQ            = 1
	HEARTBEAT_REPLY          = 2
	CONFIG_NODE_REQ          = 3
	CONFIG_NODE_REPLY        = 4
	UPDATE_GLASS_COLOR_REQ   = 5
	UPDATE_GLASS_COLOR_REPLY = 6
	GET_GLASS_STATE_REQ      = 9
	GET_GLASS_STATE_REPLY    = 10

	REPLY_CMD_START_INDEX           = 2
	REPLY_GATEWAY_ID_START_INDEX    = 3
	GATEWAY_ID_LEN                  = 6
	REPLY_NODE_ID_START_INDEX       = REPLY_GATEWAY_ID_START_INDEX + GATEWAY_ID_LEN // 9
	NODE_ID_LEN                     = 4
	HEARTBEAT_COLOR_POS_START_INDEX = REPLY_NODE_ID_START_INDEX + NODE_ID_LEN // 13
)

func findNodeIdInNodeStateList(nodeId string) *bsp.NodeState {
	for _, nodeState := range bsp.BspConfigInstance.NodeStates {
		if nodeState.NodeId == nodeId {
			return &nodeState
		}
	}
	return nil
}

func getReportedGlassColor(msg []byte) []byte {
	return msg[REPLY_NODE_ID_START_INDEX+NODE_ID_LEN+1 : len(msg)-1]
}

func getAreaFromColorReport(r byte) int {
	return (int(r)&0xf0)>>4 - 1
}

func DumpBytes(a []byte) {
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d, %d, %02x %c\n", i, a[i], a[i], a[i])
	}
	fmt.Printf("\n")
}

func HandleNodeMsg(msg []byte, mqttCh chan interface{}) {
	util.IotLogInfo("Received message\n")
	// TODO: Need to check if the message is for this gateway
	DumpBytes(msg)
	if !node.IsChecksumCorrect(msg) {
		//Log error and return
		util.IotLogErrorStr("Checksum incorrect, discard package")
		return
	}

	// Extract byte 10 to 14 as node id from byte slice msg
	nodeId := string(msg[REPLY_NODE_ID_START_INDEX : REPLY_NODE_ID_START_INDEX+NODE_ID_LEN])
	nodeState := findNodeIdInNodeStateList(nodeId)
	if nodeState == nil {
		util.IotLogInfo(fmt.Sprintf("Received heartbeat reply for unknown node %s\n", nodeId))
		return
	}
	nodeState.LastMsgTimestamp = time.Now().Unix()

	switch msg[REPLY_CMD_START_INDEX] {
	case CONFIG_NODE_REPLY:
		// TODO: handle node config reply and check if we can already update module1 & 2 config
		util.IotLogInfo(fmt.Sprintf("Received node config reply for node %s\n", nodeId))
		// TODO: handle node config reply and check if we can already update module1 & 2 config
	case HEARTBEAT_REPLY:
		util.IotLogInfo("Received heartbeat reply\n")
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+i]) & 0xf
		}

	case UPDATE_GLASS_COLOR_REPLY:
		reportedColors := getReportedGlassColor(msg)
		util.IotLogInfo("Received update glass color reply\n")
		for i := 0; i < len(reportedColors); i++ {
			nodeState.NodeReportedColor[getAreaFromColorReport(reportedColors[i])] = int(reportedColors[i] & 0xf)
		}

	case GET_GLASS_STATE_REPLY:
		// TODO: complete the real handling
		util.IotLogInfo("Received get glass state reply\n")
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[15+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[15+i]) & 0xf
		}
	}

}
