package main

import "iot_go/pkg/lora_client"

func main() {
	// client := lora_client.NewLoraClient(8866, "192.168.1.38")
	client := lora_client.NewLoraClient(8866)
	// client.InitLora()
	client.Exit()
}
