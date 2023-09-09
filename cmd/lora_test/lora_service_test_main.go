package main

import (
	"flag"
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
			// for i := 0; i < len(data); i++ {
			// 	fmt.Println(data[i])
			// }
			receivingCh <- data
		}
	}
}

func main() {
	bsp.InitConfig()
	server := ""
	isReceive := 0
	initNeeded := 0

	flag.IntVar(&isReceive, "r", 0, "is start in receive mode")
	flag.StringVar(&server, "d", "192.168.1.57", "server name")
	flag.IntVar(&initNeeded, "i", 0, "Is init needed")
	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()

	// 打印
	fmt.Printf("server=%v isReceive=%v\n", server, isReceive)

	if isReceive == 1 {
		client := lora_client.NewLoraClient(8867, server)
		client.ToggleDebug()
		if initNeeded == 1 {
			client.InitLora(bsp.BspConfigInstance.InitMsgContent.Module0)
		}
		// ch := make(chan []byte)
		// go ReceiveLoop(ch, client)
		// for {
		// 	data := <-ch
		// 	for i := 0; i < len(data); i++ {
		// 		print(data[i], "\n")
		// 	}
		// }
	} else {
		client := lora_client.NewLoraClient(8866, server)
		// client.ToggleDebug()
		if initNeeded == 1 {
			client.InitLora(bsp.BspConfigInstance.InitMsgContent.Module1)
		}
		b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
		// for {
		client.Send(b)
		// 	time.Sleep(20 * time.Second)
		// }
	}

	// client.Exit()
}
