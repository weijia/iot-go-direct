package lora_module

import (
	"context"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"time"
)

var module0SendingToNodeCh = make(chan node.NodeMsgReq, 10)
var module1SendingToNodeCh = make(chan node.NodeMsgReq, 10)
var module2SendingToNodeCh = make(chan node.NodeMsgReq, 10)

var module0ReceivingCh = make(chan []byte, 10)
var module1ReceivingCh = make(chan []byte, 10)
var module2ReceivingCh = make(chan []byte, 10)

var Module0 = LoraModule{
	SendingToNodeCh: &module0SendingToNodeCh,
	ReceivingCh:     &module0ReceivingCh,
	// LoraClient:  bsp.GetModule0Client(), // client only available after board init
	LoraIndex: 0,
}
var Module1 = LoraModule{
	SendingToNodeCh: &module1SendingToNodeCh,
	ReceivingCh:     &module1ReceivingCh,
	// LoraClient:  bsp.GetModule1Client(),
	LoraIndex: 1,
}
var Module2 = LoraModule{
	SendingToNodeCh: &module2SendingToNodeCh,
	ReceivingCh:     &module2ReceivingCh,
	// LoraClient:  bsp.GetModule2Client(),
	LoraIndex: 2,
}

var ModuleList = [3]LoraModule{
	Module0, Module1, Module2,
}

type LoraModule struct {
	SendingToNodeCh *chan node.NodeMsgReq
	ReceivingCh     *chan []byte
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
		case m := <-*loraModule.SendingToNodeCh:
			// cmd := m.Data[node.REPLY_CMD_START_INDEX]
			// if cmd == node.CONFIG_NODE_REQ {
			// util.IotLog("Received send to node req: %v using %s", m, loraModule.LoraClient.Address)
			loraModule.LoraClient.Send(m.Data)
			isTimeout, msg := loraModule.IsReplyTimeout(node.GetNodeIdStr(m.Data), util.LEVEL1_WAIT_FOR_LORA_SERVICE_PUSH_DATA_TIMEOUT_SECONDS)
			// util.IotLog("Bottom level wait for reply got: %v, %v", isTimeout, msg)
			reply := node.NodeMsgReply{
				Data:      msg,
				IsTimeout: isTimeout,
			}
			*m.ReplyCh <- reply
			// }
		}
		// util.IotLogInfo("Lora Module loop running")
	}
}

func (loraModule LoraModule) IsReplyTimeout(nodeIdStr string, timeoutSeconds int) (bool, []byte) {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	for {
		select {
		case <-eventTimer.C:
			// util.IotLog("Is reply timeout return due to timeout")
			return true, nil
		case nodeReply := <-*loraModule.ReceivingCh:
			// util.IotLog("Received node msg from ch %p", loraModule.ReceivingCh)
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
	*loraModule.SendingToNodeCh <- msgReq
}
