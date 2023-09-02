package lora_client

import (
	"fmt"
	"log"
	"net/rpc"

	"iot_go/pkg/lora_shared"
)

type LoraClient struct {
	RpcClient *rpc.Client
}

func NewLoraClient(port int) *LoraClient {
	var rpcClient *rpc.Client
	var err error
	for {
		rpcClient, err = rpc.DialHTTP("tcp", fmt.Sprintf(":%d", port))
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

func (client LoraClient) InitLora() {

	args := &lora_shared.EmptyArg{}
	var reply lora_shared.ReplyResult
	err := client.RpcClient.Call("Lora.InitLora", args, &reply)
	if err != nil {
		log.Println("Lora.InitLora error:", err)
	}
}

func (client LoraClient) Exit() {
	args := &lora_shared.EmptyArg{}
	var reply lora_shared.ReplyResult
	err := client.RpcClient.Call("Lora.Exit", args, &reply)
	if err != nil {
		log.Println("Lora.Exit error:", err)
	}
}
