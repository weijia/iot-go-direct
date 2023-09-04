package main

import (
	"fmt"
	"log"

	"os"
	"os/signal"

	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_rpc"
	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kraken-hpc/go-fork"
	"github.com/spf13/viper"
)

func init() {
	fork.RegisterFunc("StartLoraService", lora_rpc.StartLoraService)
	fork.RegisterFunc("StartLoraService1", lora_rpc.StartLoraService)
	fork.RegisterFunc("StartLoraService2", lora_rpc.StartLoraService)
	fork.Init()
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
	if err := fork.Fork("StartLoraService", "/dev/spidev1.0", 8866); err != nil {
		log.Fatalf("failed to fork: %v", err)
	}
	if err := fork.Fork("StartLoraService1", "/dev/spidev2.0", 8867); err != nil {
		log.Fatalf("failed to fork: %v", err)
	}
	if err := fork.Fork("StartLoraService2", "/dev/spidev3.0", 8868); err != nil {
		log.Fatalf("failed to fork: %v", err)
	}

	bsp.InitConfig()
	fmt.Println(viper.GetString("msg_type"))

	deviceId := bsp.BspConfigInstance.GatewayNodeID
	var topic = "device/" + deviceId + "/in"

	bsp.InitBoard()

	// MQTT 连接设置
	broker := "tcp://" + bsp.BspConfigInstance.MqttParams.MqttIP + ":" +
		fmt.Sprint(bsp.BspConfigInstance.MqttParams.MqttPort)
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(bsp.BspConfigInstance.GatewayNodeID)

	// 创建 MQTT 客户端
	client := mqtt.NewClient(opts)

	// 连接到服务器
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient := util.NewMqtt(&client, deviceId)

	// 订阅主题
	token = client.Subscribe(topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", mqttMsg.Payload(), mqttMsg.Topic())
		resp := msg.HandleMsg(mqttClient, mqttMsg.Payload())
		if resp != nil {
			mqttClient.SendToServer(resp)
		}
	})
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 发布消息
	var initMsg shared.Init
	initMsg.InitMsgContent = bsp.BspConfigInstance.InitMsgContent
	initMsg.MsgType = "init"
	mqttClient.SendToServer(initMsg)

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Disconnect(250)
	fmt.Println("Disconnected")
	bsp.GetBsp().StopAllProcess()
}
