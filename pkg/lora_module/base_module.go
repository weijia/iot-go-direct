package lora_module

import (
	"context"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"time"
)

var Module0 = LoraModule{
	SendingCh:   make(chan node.NodeMsgReq, 10),
	ReceivingCh: make(chan []byte, 10),
	// LoraClient:  bsp.GetModule0Client(), // client only available after board init
	LoraIndex:   0,
}
var Module1 = LoraModule{
	SendingCh:   make(chan node.NodeMsgReq, 10),
	ReceivingCh: make(chan []byte, 10),
	// LoraClient:  bsp.GetModule1Client(),
	LoraIndex:   1,
}
var Module2 = LoraModule{
	SendingCh:   make(chan node.NodeMsgReq, 10),
	ReceivingCh: make(chan []byte, 10),
	// LoraClient:  bsp.GetModule2Client(),
	LoraIndex:   2,
}

var ModuleList = [3]LoraModule{
	Module0, Module1, Module2,
}

type LoraModule struct {
	SendingCh   chan node.NodeMsgReq
	ReceivingCh chan []byte
	// WaitingReplyTimer *time.Timer
	LoraClient *lora_client.LoraClient
	LoraIndex  int //LoraIndex starts from 0
}

func (loraModule LoraModule) MsgLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			util.IotLogInfo("Cancel of context called, quitting LoraModule.MsgLoop")
			return
		case m := <-loraModule.SendingCh:
			// cmd := m.Data[node.REPLY_CMD_START_INDEX]
			// if cmd == node.CONFIG_NODE_REQ {
				loraModule.LoraClient.Send(m.Data)
				isTimeout, msg := loraModule.IsReplyTimeout(node.GetNodeIdStr(m.Data), 3)
				reply := node.NodeMsgReply{
					Data:      msg,
					IsTimeout: isTimeout,
				}
				*m.ReplyCh <- reply
			// }
		}
	}
}



func (loraModule LoraModule) IsReplyTimeout(nodeIdStr string, timeoutSeconds int) (bool, []byte) {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	for {
		select {
		case <-eventTimer.C:
			return true, nil
		case nodeReply := <-loraModule.ReceivingCh:
			responseNodeId := node.GetNodeIdStr(nodeReply)
			if responseNodeId == nodeIdStr {
				eventTimer.Stop()
				return false, nodeReply
			} else {
				util.IotLogErrorWithFormatStr("Unexpected node message when waiting for node reply")
			}
		}
	}
}

func (loraModule LoraModule) Send(data []byte) {
	msgReq := node.NodeMsgReq{
		Data:    data,
		ReplyCh: nil,
	}
	loraModule.SendingCh <- msgReq
}