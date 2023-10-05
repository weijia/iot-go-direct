package lora_module

import (
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"sync"
)

func (loraModule LoraModule) SendAndGetReplyOrTimeout(data []byte) (bool, []byte) {
	ch := make(chan node.NodeMsgReply)
	msgReq := node.NodeMsgReq{
		Data:    data,
		ReplyCh: &ch,
	}
	loraModule.SendingCh <- msgReq
	n := node.GetReplyOrTimeout(ch)
	if !n.IsTimeout {
		return false, n.Data
	}
	return true, make([]byte, 0)
}

func (loraModule LoraModule) UpdateGlassColorForList(
	colorUpdateParam map[string]shared.UpdateGlassColorParams, wg *sync.WaitGroup) {
	for nodeIdStr, param := range colorUpdateParam {
		c := node.GetUpdateGlassColorMsg(
			nodeIdStr, param.Color)
		isTimeout, reply := loraModule.SendAndGetReplyOrTimeout(c)
		if isTimeout || !node.IsColorUpdateSuccess(reply) {
			colorUpdateParam[nodeIdStr] = shared.UpdateGlassColorParams{
				NodeId: nodeIdStr,
				Color: node.SetColorForNodeAsInvalid(param.Color),
			}
		}
	}
	wg.Done()
}
