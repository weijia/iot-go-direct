package lora_module

import (
	"context"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"time"
)

type LoraModule struct {
	SendingCh chan []byte
	ReceivingCh chan string
	// WaitingReplyTimer *time.Timer
	LoraClient *lora_client.LoraClient
}


func (loraModule LoraModule) MsgLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			util.IotLogInfo("Cancel of context called, quitting LoraModule.MsgLoop")
		case m := <-loraModule.SendingCh:
			cmd := m[node.REPLY_CMD_START_INDEX]
			if cmd == node.CONFIG_NODE_REQ {
				loraModule.LoraClient.Send(m)
				loraModule.IsReplyTimeout(node.GetNodeIdStr(m), 3)
			}
		}
	}
}

func (loraModule LoraModule) IsReplyTimeout(nodeIdStr string, timeoutSeconds int) bool {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	for {
		select {
		case <-eventTimer.C:
			return true
		case responseNodeId := <-loraModule.ReceivingCh:
			if responseNodeId == nodeIdStr {
				eventTimer.Stop()
				return false
			} else {
				util.IotLogErrorWithFormatStr("Unexpected node message when waiting for node reply")
			}
		}
	}
}