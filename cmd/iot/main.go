package main

import (
	"fmt"

	"os"
	"os/signal"

	"iot_go/pkg/bsp"
	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

func main() {

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
		resp := msg.HandleMsg(mqttMsg.Payload())
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
}
