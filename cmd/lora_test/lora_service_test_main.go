package main

import "iot_go/pkg/lora_client"

func main() {
	client := lora_client.NewLoraClient(8866)
	client.InitLora()
	client.Exit()
}
