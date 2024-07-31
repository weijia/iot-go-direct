package lora_module

import (
	"context"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"sync"
	"time"
)

var module0SendingToNodeCh = make(chan NodeMsgReq, 10)
var module1SendingToNodeCh = make(chan NodeMsgReq, 10)
var module2SendingToNodeCh = make(chan NodeMsgReq, 10)

var module0ReceivingCh = make(chan []byte, 10)
var module1ReceivingCh = make(chan []byte, 10)
var module2ReceivingCh = make(chan []byte, 10)

var Module0 = LoraModule{
	SendingToNodeCh: &module0SendingToNodeCh,
	ReceivingCh:     &module0ReceivingCh,
	// LoraClient:  bsp.GetModule0Client(), // client only available after board init
	LoraIndex: 0,
	Mutex:     &sync.Mutex{},
}
var Module1 = LoraModule{
	SendingToNodeCh: &module1SendingToNodeCh,
	ReceivingCh:     &module1ReceivingCh,
	// LoraClient:  bsp.GetModule1Client(),
	LoraIndex: 1,
	Mutex:     &sync.Mutex{},
}
var Module2 = LoraModule{
	SendingToNodeCh: &module2SendingToNodeCh,
	ReceivingCh:     &module2ReceivingCh,
	// LoraClient:  bsp.GetModule2Client(),
	LoraIndex: 2,
	Mutex:     &sync.Mutex{},
}

var ModuleList = [3]LoraModule{
	Module0, Module1, Module2,
}

type LoraModule struct {
	SendingToNodeCh *chan NodeMsgReq
	ReceivingCh     *chan []byte
	// WaitingReplyTimer *time.Timer
	LoraClient *lora_client.LoraClient
	LoraIndex  int //LoraIndex starts from 0
	Mutex      *sync.Mutex
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
			util.IotLog("Received send to node req: %v using %s", m, loraModule.LoraClient.Address)
			loraModule.LoraClient.Send(m.Data)
			isTimeout, msg := loraModule.IsReplyTimeout(node.GetNodeIdStr(m.Data), util.LEVEL1_WAIT_FOR_LORA_SERVICE_PUSH_DATA_TIMEOUT_SECONDS)
			reply := NodeMsgReply{
				Data:      msg,
				IsTimeout: isTimeout,
			}
			util.IotLog("Bottom level wait for reply got is timeout: %v, msg: %v, send to reply ch: %p", isTimeout, msg, m.ReplyCh)
			if m.ReplyCh != nil {
				*m.ReplyCh <- reply
				close(*m.ReplyCh)
			}
			// Sleep 1 second after received a message. As other node may still working on receiving this message
			time.Sleep(time.Second * 1)
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
			util.IotLog("IsReplyTimeout return due to timeout, index: %d", loraModule.LoraIndex)
			return true, nil
		case nodeReply := <-*loraModule.ReceivingCh:
			util.IotLog("Received node msg from main msg loop (lora module ch): %p", loraModule.ReceivingCh)
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

type NodeMsgReply struct {
	Data      []byte
	IsTimeout bool
}

type NodeMsgReq struct {
	Data    []byte
	ReplyCh *chan NodeMsgReply
}

func GetReplyOrTimeout(ch *chan NodeMsgReply) NodeMsgReply {
	eventTimer := time.NewTimer(time.Duration(util.LEVEL2_NODE_MSG_REPLY_TIMEOUT_SECONDS) * time.Second)
	reply := NodeMsgReply{
		Data:      nil,
		IsTimeout: true,
	}
	select {
	case reply = <-*ch:
		util.IotLog("Got reply from base module msg loop: %v", reply)
		eventTimer.Stop()
		// util.IotLog("GetReplyOrTimeout received reply from level1, returning: %v", reply)
	case <-eventTimer.C:
		util.IotLogErrorStr("Timeout for waiting for Module to reply")
	}
	return reply
}

func (loraModule LoraModule) SendWithoutReply(data []byte) {
	msgReq := NodeMsgReq{
		Data:    data,
		ReplyCh: nil,
	}
	select {  
	case *loraModule.SendingToNodeCh <- msgReq: // 尝试发送数据到channel  
		util.IotLog("Send to node channel OK")
	default: // 如果发送操作会阻塞，则执行default分支  
		util.IotLog("Send to node channel full, discard the message")  
	}
}

func (loraModule LoraModule) SendNodeMsgWithRetryOrTimeout(msg []byte, retryCnt int, expectedMsgType int) []byte {
	for i := 0; i < retryCnt; i++ {
		ch := make(chan NodeMsgReply)
		msgReq := NodeMsgReq{
			Data:    msg,
			ReplyCh: &ch,
		}
		// util.IotLog("Before sending to sendingToNodeCh")
		*loraModule.SendingToNodeCh <- msgReq
		// util.IotLog("After sending to sendingToNodeCh")
		n := GetReplyOrTimeout(&ch)
		if !n.IsTimeout {
			if n.Data[node.REPLY_CMD_START_INDEX] == byte(expectedMsgType) {
				// util.IotLog("Is not timeout, returning data: %v", n.Data)
				return n.Data
			} else {
				util.IotLog("Received reply: %v, cmd: %d, expected: %d", n.Data, n.Data[node.REPLY_CMD_START_INDEX], expectedMsgType)
			}
		} else {
			util.IotLog("Timeout, retry sending msg: %v", msg)
		}
	}
	return nil
}
