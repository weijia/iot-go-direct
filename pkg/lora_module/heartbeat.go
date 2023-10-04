package lora_module

import (
	"context"
	"sync"

	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
)

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
		reply := node.SendHeartbeatForNode(nodeIdStr, &loraModule.SendingCh)
		if reply.IsTimeout {
			// Retry twice as requested in node msg document
			reply = node.SendHeartbeatForNode(nodeIdStr, &loraModule.SendingCh)
			if reply.IsTimeout {
				node.SendNodeInitForNode(nodeIdStr, module, &loraModule.SendingCh)
			}
		}
	}
	if wg != nil {
		wg.Done()
	}
}
