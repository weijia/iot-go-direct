package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"

	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_rpc"
	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"

	"github.com/kraken-hpc/go-fork"
)

func init() {
	fork.RegisterFunc("StartLoraService", lora_rpc.StartLoraServiceInBackgrand)
	fork.RegisterFunc("StartLoraService1", lora_rpc.StartLoraServiceInBackgrand)
	fork.RegisterFunc("StartLoraService2", lora_rpc.StartLoraServiceInBackgrand)
	fork.Init()
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
	var loraHost string
	flag.StringVar(&loraHost, "s", "127.0.0.1", "Lora service server")
	// loraServiceIp := loraHost
	loraServiceIp := "192.168.1.46"
	// gatewayPushMsgIp := "127.0.0.1"
	gatewayPushMsgIp := "192.168.1.30"

	if err := fork.Fork("StartLoraService", "/dev/spidev1.0", 8866, gatewayPushMsgIp); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("failed to fork: %v", err))
	}
	if err := fork.Fork("StartLoraService1", "/dev/spidev2.0", 8867, gatewayPushMsgIp); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("failed to fork: %v", err))
	}
	if err := fork.Fork("StartLoraService2", "/dev/spidev3.0", 8868, gatewayPushMsgIp); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("failed to fork: %v", err))
	}

	bsp.InitConfig()
	// fmt.Println(viper.GetString("msg_type"))

	bsp.InitBoard(loraServiceIp)

	quit := make(chan bool)

	nodeMsgCh := make(chan []byte, 10)

	go lora_rpc.StartLoraReceiverRpc(&nodeMsgCh, 8869)
	go msg.StartMqttMsgLoop(nodeMsgCh, quit)

	// go node.SendNodeHeartbeatInLoop()
	// node.SendNodeInitAfterStartup()

	// 发布消息
	var initMsg shared.Init
	initMsg.InitMsgContent = bsp.BspConfigInstance.InitMsgContent
	initMsg.MsgType = "init"
	msg.MqttPublishCh <- initMsg

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	quit <- true

	bsp.GetBsp().StopAllProcess()
}
