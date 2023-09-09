package lora_client

import (
	"fmt"
	"log"
	"net/rpc"

	"iot_go/pkg/lora_shared"
	"iot_go/pkg/shared"
)

type LoraClient struct {
	RpcClient *rpc.Client
}

func NewLoraClient(port int, host ...string) *LoraClient {
	var rpcClient *rpc.Client
	var err error
	realHost := "127.0.0.1"
	if len(host) > 0 {
		realHost = host[0]
	}
	for {
		rpcClient, err = rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", realHost, port))
		if err != nil {
			log.Println("dialing:", err)
		} else {
			break
		}
	}
	instance := new(LoraClient)
	instance.RpcClient = rpcClient
	return instance
}

func (client LoraClient) InitLora(module shared.Module) {
	var reply lora_shared.ReplyResult
	err := client.RpcClient.Call("LoraRpc.InitLora", module, &reply)
	if err != nil {
		log.Println("Lora.InitLora error:", err)
	}
}

func (client LoraClient) Exit() {
	args := &lora_shared.EmptyArg{}
	var reply lora_shared.ReplyResult
	err := client.RpcClient.Call("LoraRpc.Exit", args, &reply)
	if err != nil {
		log.Println("LoraRpc.Exit error:", err)
	}
}

func (client LoraClient) Send(data []byte) {
	args := &lora_shared.LoraData{Data: data}
	var reply lora_shared.ReplyResult
	err := client.RpcClient.Call("LoraRpc.Send", args, &reply)
	if err != nil {
		log.Println("LoraRpc.Send error:", err)
	}
}

func (client LoraClient) Receive() []byte {
	args := &lora_shared.EmptyArg{}
	var reply lora_shared.LoraData
	err := client.RpcClient.Call("LoraRpc.Receive", args, &reply)
	if err != nil {
		log.Println("LoraRpc.Receive error:", err)
	}
	for i := 0; i < len(reply.Data); i++ {
		log.Printf("0x%x", reply.Data[i])
	}
	log.Printf("\n")
	return reply.Data
}

func (client LoraClient) ToggleDebug() {
	args := &lora_shared.EmptyArg{}
	var reply lora_shared.EmptyArg
	client.RpcClient.Call("LoraRpc.ToggleDebug", args, &reply)
}

func (client LoraClient) OnReceive(data []byte) {
	args := &lora_shared.LoraData{Data: data}
	var reply lora_shared.ReplyResult
	err := client.RpcClient.Call("LoraReceiverRpc.OnReceive", args, &reply)
	if err != nil {
		log.Println("LoraReceiverRpc.OnReceive error:", err)
	}
}
