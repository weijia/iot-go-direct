package lora_module

import (
	"context"
	"sync"

	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)

// Will be set by top level msg loop
var TimeoutNodeIdCh *chan string

func (loraModule LoraModule) SendHeartbeatForList(ctx context.Context,
	nodeList []string, wg *sync.WaitGroup) {

	module := bsp.BspConfigInstance.Module0

	if loraModule.LoraIndex == 1 {
		module = bsp.BspConfigInstance.Module1
	} else {
		if loraModule.LoraIndex == 2 {
			module = bsp.BspConfigInstance.Module2
		}
	}
	for _, nodeIdStr := range nodeList {
		util.IotLog("Sending heartbeat on %d for %s", loraModule.LoraIndex, nodeIdStr)
		reply := loraModule.SendNodeMsgWithRetryOrTimeout(node.GetHeartbeatMsg(nodeIdStr), util.HEARTBEAT_RETRY_CNT, node.HEARTBEAT_REPLY)
		// util.IotLog("SendHeartbeatForNode returned timeout & data: %v", reply)
		if reply == nil {
			if TimeoutNodeIdCh != nil {
				*TimeoutNodeIdCh <- nodeIdStr
			}
			util.IotLogErrorWithFormatStr("Heartbeat for %s no reply, will set node "+
				"status as offline and send node init for it on public freq", nodeIdStr)
			Module0.SendNodeMsgWithRetryOrTimeout(node.GetNodeInitMsg(nodeIdStr, module), util.NODE_INIT_RETRY_CNT, node.HEARTBEAT_REPLY)
		}
	}
	if wg != nil {
		wg.Done()
	}
}
