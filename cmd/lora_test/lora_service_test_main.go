package main

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"time"
)

var isQuit = false

func ReceiveLoop(receivingCh chan<- []byte, client *lora_client.LoraClient) {
	for {
		if isQuit {
			break
		}
		data := client.Receive()
		if len(data) == 0 {
			fmt.Println("No data received, wait for 10 seconds")
			time.Sleep(10 * time.Second)
		} else {
			for i := 0; i < len(data); i++ {
				print(data[i], "\n")
			}
			receivingCh <- data
		}
	}
}

func main() {
	bsp.InitConfig()
	client := lora_client.NewLoraClient(8867, "192.168.1.54")
	// client := lora_client.NewLoraClient(8866, "192.168.1.54")
	// client := lora_client.NewLoraClient(8866)
	// client.InitLora(bsp.BspConfigInstance.InitMsgContent.Module0)
	// client.InitLora(bsp.BspConfigInstance.InitMsgContent.Module1)
	// ch := make(chan []byte)
	// go ReceiveLoop(ch, client)
	// data := <-ch
	// for i := 0; i < len(data); i++ {
	// 	print(data[i], "\n")
	// }
	// isQuit = true
	b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	client.Send(b)
	// client.Exit()
}
