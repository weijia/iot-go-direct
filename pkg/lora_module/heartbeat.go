package lora_module

import (
	"context"
	"sync"

	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)
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
		util.IotLog("Sending heartbeat for %s", nodeIdStr)
		reply := node.SendHeartbeatForNode(nodeIdStr, &loraModule.SendingCh)
		if reply.IsTimeout {
			// Retry twice as requested in node msg document
			util.IotLog("Sending heartbeat for %s 2nd time", nodeIdStr)
			reply = node.SendHeartbeatForNode(nodeIdStr, &loraModule.SendingCh)
			if reply.IsTimeout {
				if TimeoutNodeIdCh != nil {
					*TimeoutNodeIdCh <- nodeIdStr
				}
				util.IotLogErrorWithFormatStr("Heartbeat for %s no reply, will set node " + 
					"status as offline and send node init for it on public freq", nodeIdStr)
				node.SendNodeInitForNode(nodeIdStr, module, &Module0.SendingCh)
			}
		}
	}
	if wg != nil {
		wg.Done()
	}
}
