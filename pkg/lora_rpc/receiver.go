package lora_rpc

import (
	"fmt"
	"iot_go/pkg/lora_shared"
	"iot_go/pkg/util"
	"net/http"
	"net/rpc"
)

type LoraReceiverRpc struct {
}

var recvChannel *chan lora_shared.LoraData

func (loraReceiverRpc LoraReceiverRpc) OnReceive(argType lora_shared.LoraData, reply *lora_shared.ReplyResult) error {
	// util.IotLog("RPC:  called")
	select {  
		case *recvChannel <- argType:  
			// fmt.Println("Sent successfully")  
		default:  
			// 如果channel已满，打印消息  
			util.IotLogErrorStr("Channel full, could not send without blocking")
	}
	// log.Println("RPC: after put to channel")
	util.IotLog("---------------Received lora msg from hardware %d: % X", argType.ModuleIndex, argType.Data)
	reply.Result = 0
	return nil
}

var ReceivingServer *http.Server

func StartLoraReceiverRpc(recvCh *chan lora_shared.LoraData, port int) {
	recvChannel = recvCh
	loraRpc := LoraReceiverRpc{}
	// The loraRpc object will be copied to RPC procedure instead of sending the original object
	// So the data changed after Register may be discarded
	rpc.Register(loraRpc)
	rpc.HandleHTTP()
	ReceivingServer = &http.Server{Addr: fmt.Sprintf(":%d", port)}
	if err := ReceivingServer.ListenAndServe(); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("serve error: %v", err))
	}
}
