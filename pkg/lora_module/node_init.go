package lora_module

import (
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"sync"
)

func (loraModule LoraModule) SendNodeInitForList(
	// ctx context.Context,
	nodeParam map[string]shared.Module,
	wg *sync.WaitGroup) {
	for nodeIdStr, param := range nodeParam {
		// loraModule.Mutex.Lock()
		loraModule.SendNodeMsgWithRetryOrTimeout(node.GetNodeInitMsg(nodeIdStr, param), util.NODE_INIT_RETRY_CNT, node.CONFIG_NODE_REPLY)
		// loraModule.Mutex.Unlock()
	}
	wg.Done()
}
