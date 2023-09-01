package main

import (
	"iot_go/pkg/lora"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":8866")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &lora.EmptyArg{}
	var reply lora.ReplyResult
	err = client.Call("Lora.InitLora", args, &reply)
	if err != nil {
		log.Fatal("Lora.InitLora error:", err)
	}
	err = client.Call("Lora.Exit", args, &reply)
	if err != nil {
		log.Fatal("Lora.Exit error:", err)
	}
}
