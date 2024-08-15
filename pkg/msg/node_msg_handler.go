package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/lora_shared"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"strings"
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

// func getAreaFromColorReport(r byte) int {
// 	return (int(r)&0xf0)>>4 - 1
// }

func HandleNodeMsg(msgRaw lora_shared.LoraData, MqttPublishCh chan interface{}) {
	// Application's data change will only happen in this routine to ensure no
	// concurrent data changes in different routines
	util.IotDebugPrintf("Handle Node Msg: Received node message from module index: %d, %v", msgRaw.ModuleIndex, msgRaw.Data)
	msg := msgRaw.Data
	// DumpBytes(msg)
	if !node.IsChecksumCorrect(msg) {
		//Log error and return
		util.IotLogErrorWithFormatStr("Handle Node Msg: Checksum incorrect, discard package: %v", msg)
		return
	}
	gatewayId := fmt.Sprintf("%x", msg[REPLY_GATEWAY_ID_START_INDEX:REPLY_GATEWAY_ID_START_INDEX+GATEWAY_ID_LEN])

	// Need to check if the message is for this gateway
	if !strings.EqualFold(gatewayId, bsp.BspConfigInstance.GatewayNodeId) {
		util.IotDebug(fmt.Sprintf("Handle Node Msg: Received gateway Id: %s is not for this gateway: %s", gatewayId, bsp.BspConfigInstance.GatewayNodeId))
		return
	}

	// Extract byte 10 to 14 as node id from byte slice msg
	nodeId := node.GetNodeIdStr(msg)
	nodeState := bsp.GetOrCreateNodeState(nodeId)
	if nodeState == nil {
		util.IotLogErrorWithFormatStr("Received msg but can not create/get node state: %s", nodeId)
		return
	}
	nodeState.LastMsgTimestamp = time.Now().Unix()
	// util.IotLog("LastMsgTimestamp: %d", nodeState.LastMsgTimestamp)

	switch msg[REPLY_CMD_START_INDEX] {
	case CONFIG_NODE_REPLY:
		// TODO: handle node config reply and check if we can already update module1 & 2 config
		util.IotDebug(fmt.Sprintf("Handle Node Msg: Received node config reply for node %s", nodeId))
		node.UpdateNodeStateForInitReply(msg)
		util.SendBytesMsgWithoutBlocking(msg, lora_module.ModuleList[msgRaw.ModuleIndex].ReceivingCh,
			fmt.Sprintf("Send received node reply %v to module %d failed", msg, msgRaw.ModuleIndex))

	case UNKNOWN_NODE_REPLY:
		util.IotDebug(fmt.Sprintf("Handle Node Msg: Received unknown node config reply for node %s", nodeId))
		node.UpdateNodeStateForInitReply(msg)
		if !bsp.IsInNodeList1(nodeId) && !bsp.IsInNodeList2(nodeId) {
			var unknownNodeList []string
			unknownNodeList = append(unknownNodeList, nodeId)
			reply := UnknownNode{
				MsgType:         "unknown_node",
				UnknownNodeList: unknownNodeList,
			}
			MqttPublishCh <- reply
		}

	case HEARTBEAT_REPLY:
		util.IotDebug("Handle Node Msg: Received is heartbeat reply")
		// util.IotLog("Before update according to heartbeat: %v", nodeState.NodeReportedColor)
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+i]) & 0xf
		}
		nodeState.CompletionStatus = int(msg[HEARTBEAT_COLOR_POS_START_INDEX+4])
		nodeState.RSSI = msgRaw.RSSI
		nodeState.SNR = msgRaw.SNR
		nodeState.IsOffline = false
		// util.IotLog("Handle Node Msg: After update according to heartbeat: %v", nodeState.NodeReportedColor)
		// util.IotLog("module: %d, ch: %p", msgRaw.ModuleIndex, lora_module.ModuleList[msgRaw.ModuleIndex].ReceivingCh)
		util.SendBytesMsgWithoutBlocking(msg, lora_module.ModuleList[msgRaw.ModuleIndex].ReceivingCh,
			fmt.Sprintf("Handle Node Msg: Send received heartbeat reply %v to module %d failed", msg, msgRaw.ModuleIndex))

	case UPDATE_GLASS_COLOR_REPLY:
		util.IotDebug("Handle Node Msg: Received is glass color update reply")
		util.IotDebugPrintf("Handle Node Msg: module: %d, ch: %p", msgRaw.ModuleIndex, lora_module.ModuleList[msgRaw.ModuleIndex].ReceivingCh)
		util.SendBytesMsgWithoutBlocking(msg, lora_module.ModuleList[msgRaw.ModuleIndex].ReceivingCh,
			fmt.Sprintf("Handle Node Msg: Send received update glass color reply %v to module %d failed", msg, msgRaw.ModuleIndex))

	case GET_GLASS_STATE_REPLY:
		// TODO: complete the real handling
		util.IotDebug("Handle Node Msg: Received is get glass state reply")
		for i := 0; i < 4; i++ {
			nodeState.NodeReportedColor[i*2] = int(msg[node.REPLY_PAYLOAD_START_INDEX+i]) & 0xf0 >> 4
			nodeState.NodeReportedColor[i*2+1] = int(msg[node.REPLY_PAYLOAD_START_INDEX+i]) & 0xf
		}
		util.SendBytesMsgWithoutBlocking(msg, lora_module.ModuleList[msgRaw.ModuleIndex].ReceivingCh,
			fmt.Sprintf("Handle Node Msg: Send received glass state reply %v to module %d failed", msg, msgRaw.ModuleIndex))
	default:
		util.IotLogErrorStr("Handle Node Msg: Not handled msg, maybe it is sent from gateway")
	}

}
