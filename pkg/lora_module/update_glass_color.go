package lora_module

import (
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"sync"
)

func (loraModule LoraModule) UpdateGlassColorForList(
	colorUpdateParam map[string]shared.UpdateGlassColorParams, wg *sync.WaitGroup) {
	loraModule.Mutex.Lock()
	for nodeIdStr, param := range colorUpdateParam {
		c := node.GetUpdateGlassColorMsg(
			nodeIdStr, param.Color)
		reply := loraModule.SendNodeMsgWithRetryOrTimeout(c, util.UPDATE_GLASS_COLOR_RETRY_CNT, node.UPDATE_GLASS_COLOR_REPLY)
		if reply == nil || !node.IsColorUpdateSuccess(reply) {
			colorUpdateParam[nodeIdStr] = shared.UpdateGlassColorParams{
				NodeId: nodeIdStr,
				Color:  node.SetColorForNodeAsInvalid(param.Color),
			}
		}
	}
	loraModule.Mutex.Unlock()
	wg.Done()
}
