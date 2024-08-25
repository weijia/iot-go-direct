package lora_module

import (
	"context"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"time"
)

var module0SendingToNodeCh = make(chan NodeMsgReq, 10)
var module1SendingToNodeCh = make(chan NodeMsgReq, 10)
var module2SendingToNodeCh = make(chan NodeMsgReq, 10)

var module0ReceivingCh = make(chan []byte, 10)
var module1ReceivingCh = make(chan []byte, 10)
var module2ReceivingCh = make(chan []byte, 10)

var module1NodeListCh = make(chan []string, 2)
var module2NodeListCh = make(chan []string, 2)

var Module0 = LoraModule{
	SendingToNodeCh: &module0SendingToNodeCh,
	ReceivingCh:     &module0ReceivingCh,
	// LoraClient:  bsp.GetModule0Client(), // client only available after board init
	LoraIndex: 0,
	// Mutex:     &sync.Mutex{},
}
var Module1 = LoraModule{
	SendingToNodeCh: &module1SendingToNodeCh,
	ReceivingCh:     &module1ReceivingCh,
	// LoraClient:  bsp.GetModule1Client(),
	LoraIndex: 1,
	// Mutex:     &sync.Mutex{},
	NodeListCh: &module1NodeListCh,
}

var Module2 = LoraModule{
	SendingToNodeCh: &module2SendingToNodeCh,
	ReceivingCh:     &module2ReceivingCh,
	// LoraClient:  bsp.GetModule2Client(),
	LoraIndex: 2,
	// Mutex:     &sync.Mutex{},
	NodeListCh: &module2NodeListCh,
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
	// Mutex      *sync.Mutex
	// nodeList []string
	NodeListCh *chan []string
}

func (loraModule LoraModule) MsgLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			util.IotLogInfo("NodeMsgLoop: Cancel of context called, quitting LoraModule.MsgLoop")
			return
		case nodeReply := <-*loraModule.ReceivingCh:
			util.IotLogErrorWithFormatStr("NodeMsgLoop: unexpected node message when no request message sent: %v", nodeReply)
		case m := <-*loraModule.SendingToNodeCh:
			util.IotDebugPrintf("NodeMsgLoop: Received send to node req: %v using %s", m, loraModule.LoraClient.Address)
			loraModule.LoraClient.Send(m.Data)
			isTimeout, msg := loraModule.IsReplyTimeoutWillBlock(node.GetNodeIdStr(m.Data), util.LEVEL1_WAIT_FOR_LORA_SERVICE_PUSH_DATA_TIMEOUT_SECONDS)
			reply := NodeMsgReply{
				Data:      msg,
				IsTimeout: isTimeout,
			}
			util.IotDebugPrintf("NodeMsgLoop: Bottom level wait for reply got is timeout: %v, msg: %v, send to reply ch: %p", isTimeout, msg, m.ReplyCh)
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

func (loraModule LoraModule) IsReplyTimeoutWillBlock(nodeIdStr string, timeoutSeconds int) (bool, []byte) {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	defer eventTimer.Stop()
	for {
		select {
		case <-eventTimer.C:
			util.IotLogErrorWithFormatStr("NodeMsgLoop: IsReplyTimeoutWillBlock return due to timeout, index: %d", loraModule.LoraIndex)
			return true, nil
		case nodeReply := <-*loraModule.ReceivingCh:
			util.IotDebugPrintf("NodeMsgLoop: Received node msg from main msg loop (lora module ch): %p", loraModule.ReceivingCh)
			responseNodeId := node.GetNodeIdStr(nodeReply)
			if responseNodeId == nodeIdStr {
				// eventTimer.Stop() used defer above to do this already
				return false, nodeReply
			} else {
				util.IotLogErrorWithFormatStr("NodeMsgLoop: IsReplyTimeoutWillBlock Unexpected node message when waiting for node reply, response node: %s, expected node: %s", responseNodeId, nodeIdStr)
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
	defer eventTimer.Stop()
	reply := NodeMsgReply{
		Data:      nil,
		IsTimeout: true,
	}
	select {
	case reply = <-*ch:
		util.IotDebugPrintf("GetReplyOrTimeout: Got reply from base module msg loop (data, is timeout): %v", reply)
		// eventTimer.Stop() // already use defer above to do this
		// util.IotLog("GetReplyOrTimeout received reply from level1, returning: %v", reply)
	case <-eventTimer.C:
		util.IotLogErrorStr("GetReplyOrTimeout: Timeout for waiting for Module to reply")
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
		// util.IotLog("Send to node channel OK")
	default: // 如果发送操作会阻塞，则执行default分支
		util.IotLogErrorStr("SendWithoutReply: Send to node channel full, discard the message")
	}
}

func (loraModule LoraModule) SendNodeMsgWithRetryOrTimeoutWillBlock(msg []byte, retryCnt int, expectedMsgType int) []byte {
	for i := 0; i < retryCnt; i++ {
		ch := make(chan NodeMsgReply)
		msgReq := NodeMsgReq{
			Data:    msg,
			ReplyCh: &ch,
		}
		// util.IotLog("Before sending to sendingToNodeCh")
		select {
		case *loraModule.SendingToNodeCh <- msgReq:
		default:
			// 如果channel已满，打印消息
			util.IotLogErrorStr("SendNodeMsgWithRetryOrTimeoutWillBlock: Channel full, could not send without blocking")
			return nil
		}

		// util.IotLog("After sending to sendingToNodeCh")
		n := GetReplyOrTimeout(&ch)
		if !n.IsTimeout {
			if n.Data[node.REPLY_CMD_START_INDEX] == byte(expectedMsgType) {
				// util.IotLog("Is not timeout, returning data: %v", n.Data)
				return n.Data
			} else {
				util.IotLogErrorWithFormatStr("SendNodeMsgWithRetryOrTimeoutWillBlock: Received reply: %v, cmd: %d, expected: %d", n.Data, n.Data[node.REPLY_CMD_START_INDEX], expectedMsgType)
			}
		} else {
			util.IotLogErrorWithFormatStr("SendNodeMsgWithRetryOrTimeoutWillBlock: Timeout, retry sending msg: %v", msg)
		}
	}
	return nil
}
