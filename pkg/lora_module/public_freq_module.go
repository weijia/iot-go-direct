package lora_module

import (
	"context"
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)

var PublicFreqNodeMsgCh = make(chan []byte, 10)

func SendMsgOnPublicFreq(m []byte) {
	util.SendMsgWithoutBlocking(m, PublicFreqNodeMsgCh, "Timeout for sending msg")
}

func PublicFreqModuleMsgLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			util.IotLogInfo("Cancel of context called, quitting PublicFreqModuleMsgLoop")
		case m := <-PublicFreqNodeMsgCh:
			cmd := m[node.REPLY_CMD_START_INDEX]
			if cmd == node.CONFIG_NODE_REQ {
				bsp.GetModule0Client().Send(m)
			}
		}
	}
}