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
	UPDATE_GLASS_COLOR_REQ   = 5
	UPDATE_GLASS_COLOR_REPLY = 6
	NODE_ID_START_POS        = 10
	NODE_ID_LEN              = 4
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
	return msg[NODE_ID_START_POS+NODE_ID_LEN+1 : len(msg)-1]
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
	DumpBytes(msg)
	if !node.IsChecksumCorrect(msg) {
		//Log error and return
		util.IotLogErrorStr("Checksum incorrect, discard package")
		return
	}

	// Extract byte 10 to 14 as node id from byte slice msg
	nodeId := string(msg[NODE_ID_START_POS : NODE_ID_START_POS+NODE_ID_LEN])
	nodeState := findNodeIdInNodeStateList(nodeId)
	if nodeState == nil {
		util.IotLogInfo(fmt.Sprintf("Received heartbeat reply for unknown node %s\n", nodeId))
		return
	}
	nodeState.LastMsgTimestamp = time.Now().Unix()

	switch msg[2] {
	case HEARTBEAT_REPLY:
		util.IotLogInfo("Received heartbeat reply\n")
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[15+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[15+i]) & 0xf
		}

	case UPDATE_GLASS_COLOR_REPLY:
		reportedColors := getReportedGlassColor(msg)
		util.IotLogInfo("Received update glass color reply\n")
		for i := 0; i < len(reportedColors); i++ {
			nodeState.NodeReportedColor[getAreaFromColorReport(reportedColors[i])] = int(reportedColors[i] & 0xf)
		}
	}
}
