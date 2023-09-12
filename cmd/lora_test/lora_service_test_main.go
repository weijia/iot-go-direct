package main

import (
	"flag"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
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
	port := 8866

	flag.IntVar(&isReceive, "r", 0, "is start in receive mode")
	flag.StringVar(&server, "d", "192.168.1.57", "server name")
	flag.IntVar(&initNeeded, "i", 0, "Is init needed")
	flag.IntVar(&port, "p", 0, "Server port")
	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()

	// 打印
	fmt.Printf("server=%v isReceive=%v\n", server, isReceive)

	if isReceive == 1 {
		client := lora_client.NewLoraClient(port, server)
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
		client := lora_client.NewLoraClient(port, server)
		// client.ToggleDebug()
		if initNeeded == 1 {
			client.InitLora(bsp.BspConfigInstance.InitMsgContent.Module0)
		}
		// b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
		// b := node.GetNodeInitMsg(
		// 	bsp.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		// 	bsp.DecodeId("01020304"), bsp.BspConfigInstance.InitMsgContent.Module1)
		// h := node.GetHeartBeatMsg(
		// 	bsp.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		// 	bsp.DecodeId("01020304"))
		// c := node.GetUpdateGlassColorMsg(
		// 	bsp.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		// 	bsp.DecodeId("01020304"), "122331")
		// nodeGroup := []([]byte){
		// 	bsp.DecodeId("01020304"),
		// 	bsp.DecodeId("01020305"),
		// }
		// g := node.GetGroupUpdateGlassColorMsg(
		// 	bsp.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		// 	nodeGroup, "122331")
		// b1 := node.GetBroadcastUpdateGlassColorMsg(
		// 	bsp.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		// 	"122331")
		b2 := node.GetRetrieveColorMsg(
			bsp.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
			bsp.DecodeId("01020304"))
		// for {
		// client.Send(b)
		// client.Send(h)
		// client.Send(c)
		// client.Send(g)
		// client.Send(b1)
		client.Send(b2)
		// 	time.Sleep(20 * time.Second)
		// }
	}

	// client.Exit()
}
