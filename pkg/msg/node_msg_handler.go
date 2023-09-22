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
	UNKNOWN_NODE_REPLY       = 40
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
	util.IotLogInfo("Received node message")
	// TODO: Need to check if the message is for this gateway
	DumpBytes(msg)
	if !node.IsChecksumCorrect(msg) {
		//Log error and return
		util.IotLogErrorStr("Checksum incorrect, discard package")
		return
	}

	// Extract byte 10 to 14 as node id from byte slice msg
	nodeId := fmt.Sprintf("%x", msg[REPLY_NODE_ID_START_INDEX:REPLY_NODE_ID_START_INDEX+NODE_ID_LEN])
	nodeState := bsp.GetNodeState(nodeId)
	if nodeState == nil {
		util.IotLogInfo(fmt.Sprintf("Received reply for unknown node: %s", nodeId))
		return
	}
	nodeState.LastMsgTimestamp = time.Now().Unix()

	switch msg[REPLY_CMD_START_INDEX] {
	case CONFIG_NODE_REPLY:
		// TODO: handle node config reply and check if we can already update module1 & 2 config
		util.IotLogInfo(fmt.Sprintf("Received node config reply for node %s", nodeId))
		node.HandleNodeInitReply(msg)

	case UNKNOWN_NODE_REPLY:
		// TODO: handle node config reply and check if we can already update module1 & 2 config
		util.IotLogInfo(fmt.Sprintf("Received unknown node config reply for node %s", nodeId))
		node.HandleNodeInitReply(msg)

	case HEARTBEAT_REPLY:
		util.IotLogInfo("Received heartbeat reply")
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+i]) & 0xf
		}

		util.SendRepliedNodeIdWithoutBlocking(nodeId, node.HeartbeatCh, 10)

	case UPDATE_GLASS_COLOR_REPLY:
		util.IotLogInfo("Received update glass color reply")
		// node.HandleNodeColorUpdateReply(msg)
		for pendingUpdateIndex := len(node.StatesForPendingColorUpdate) - 1; pendingUpdateIndex >= 0; pendingUpdateIndex-- {
			state := node.StatesForPendingColorUpdate[pendingUpdateIndex]
			for _, singleNodeState := range state.Reply.Status {
				if singleNodeState.NodeId == nodeId {
					// "params":[
					// {
					// 	"node_id":"FD000001",
					// 	"color":"1223"
					// },
					if int(msg[node.REPLY_PAYLOAD_START_INDEX]) == node.UPDATE_COLOR_RESULT_OK {
						for i := 0; i < len(singleNodeState.Color)/2; i++ {
							nodeState.NodeReportedColor[util.GetGlassAreaFromStr(singleNodeState.Color[i*2])] =
								int(singleNodeState.Color[i*2+1])
						}
					} else {
						singleNodeState.Color = node.SetColorForNodeAsInvalid(singleNodeState.Color)
					}
					// Remove pending node, will only remove one, so we can remove when using range
					for pendingNodeIndex, pendingNodeId := range state.PendingNodes {
						if pendingNodeId == nodeId {
							state.PendingNodes = append(state.PendingNodes[:pendingNodeIndex],
								state.PendingNodes[pendingNodeIndex+1:]...)
							break
						}
					}
				}
			}
			if len(state.PendingNodes) == 0 {
				// Remove the state if all pending nodes sent reply
				node.StatesForPendingColorUpdate =
					append(node.StatesForPendingColorUpdate[:pendingUpdateIndex],
						node.StatesForPendingColorUpdate[pendingUpdateIndex+1:]...)
				mqttCh <- state.Reply
			}
		}

	case GET_GLASS_STATE_REPLY:
		// TODO: complete the real handling
		util.IotLogInfo("Received get glass state reply")
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[15+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[15+i]) & 0xf
		}
	}

}
