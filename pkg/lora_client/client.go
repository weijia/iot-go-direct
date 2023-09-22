package lora_client

import (
	"fmt"
	"log"
	"net/rpc"

	"iot_go/pkg/lora_shared"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

type LoraClient struct {
	RpcClient *rpc.Client
	Address   string
}

func (client *LoraClient) CreateRpcClientWithRetry() {
	var err error
	for {
		client.RpcClient, err = rpc.DialHTTP("tcp", client.Address)
		if err != nil {
			log.Println("dialing:", err)
		} else {
			break
		}
	}
}

func NewLoraClient(port int, host ...string) *LoraClient {
	realHost := "127.0.0.1"
	if len(host) > 0 {
		realHost = host[0]
	}
	instance := new(LoraClient)
	instance.Address = fmt.Sprintf("%s:%d", realHost, port)
	instance.CreateRpcClientWithRetry()
	return instance
}

func (client LoraClient) CallWithReconnect(serviceMethod string, args any, reply any) error {
	err := client.RpcClient.Call(serviceMethod, args, reply)
	if rpc.ErrShutdown == err {
		client.CreateRpcClientWithRetry()
		err = client.RpcClient.Call(serviceMethod, args, reply)
	}
	return err
}

func (client LoraClient) InitLora(module shared.Module) {
	var reply lora_shared.ReplyResult
	err := client.CallWithReconnect("LoraRpc.InitLora", module, &reply)
	if err != nil {
		log.Println("Lora.InitLora error:", err)
	}
}

func (client LoraClient) Exit() {
	args := &lora_shared.EmptyArg{}
	var reply lora_shared.ReplyResult
	err := client.CallWithReconnect("LoraRpc.Exit", args, &reply)
	if err != nil {
		log.Println("LoraRpc.Exit error:", err)
	}
}

func (client LoraClient) Send(data []byte) {
	args := &lora_shared.LoraData{Data: data}
	var reply lora_shared.ReplyResult
	err := client.CallWithReconnect("LoraRpc.Send", args, &reply)
	if err != nil {
		log.Println("LoraRpc.Send error:", err)
	}
}

func (client LoraClient) Receive() []byte {
	args := &lora_shared.EmptyArg{}
	var reply lora_shared.LoraData
	err := client.CallWithReconnect("LoraRpc.Receive", args, &reply)
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
	client.CallWithReconnect("LoraRpc.ToggleDebug", args, &reply)
}

func (client LoraClient) OnReceive(data []byte) {
	args := &lora_shared.LoraData{Data: data}
	var reply lora_shared.ReplyResult
	// util.IotLogInfo("Before on receive call\n")
	err := client.CallWithReconnect("LoraReceiverRpc.OnReceive", args, &reply)
	if err != nil {
		util.IotLogInfo("Error in on receive call\n")
		util.IotLogError(err)
	}
	// util.IotLogInfo("After on receive call\n")
}
