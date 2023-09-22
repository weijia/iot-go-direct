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

var recvChannel *chan []byte

func (loraReceiverRpc LoraReceiverRpc) OnReceive(argType lora_shared.LoraData, reply *lora_shared.ReplyResult) error {
	// log.Println("RPC: OnReceive called")
	*recvChannel <- argType.Data
	// log.Println("RPC: after put to channel")
	reply.Result = 0
	return nil
}

func StartLoraReceiverRpc(recvCh *chan []byte, port int) {
	recvChannel = recvCh
	rolaRpc := LoraReceiverRpc{}
	// The rolaRpc object will be copied to RPC procedure instead of sending the original object
	// So the data changed after Register may be discarded
	rpc.Register(rolaRpc)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("serve error: %v", err))
	}
}
