package lora_module

import (
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"sync"
)

func (loraModule LoraModule) SendNodeInitForList(
	// ctx context.Context,
	nodeParam map[string]shared.Module,
	wg *sync.WaitGroup) {
	for nodeIdStr, param := range nodeParam {
		node.SendNodeInitForNode(nodeIdStr, param, &loraModule.SendingCh)
	}
	wg.Done()
}
