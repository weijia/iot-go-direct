package main

import (
	"encoding/json"
	"fmt"

	"os"
	"os/signal"
	"time"

	"iot_go/pkg/bsp"
	"iot_go/pkg/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

func main() {

	bsp.InitConfig()
	fmt.Println(viper.GetString("msg_type"))

	deviceId := bsp.BspConfigInstance.InitConfig.GatewayNodeID
	var topic = "device/" + deviceId + "/in"
	var publishTopic = "device/" + deviceId + "/out"

	bsp.InitBoard()

	// MQTT 连接设置
	broker := "tcp://" + bsp.BspConfigInstance.MqttParams.MqttIP + ":" +
		fmt.Sprint(bsp.BspConfigInstance.MqttParams.MqttPort)
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(bsp.BspConfigInstance.InitConfig.GatewayNodeID)

	// 创建 MQTT 客户端
	client := mqtt.NewClient(opts)

	// 连接到服务器
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	token = client.Subscribe(topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", mqttMsg.Payload(), mqttMsg.Topic())
		resp := msg.HandleMsg(mqttMsg.Payload())
		if resp != nil {
			payload, _ := json.Marshal(resp)
			token := client.Publish(publishTopic, 0, false, payload)
			token.Wait()
			fmt.Printf("Published message: %s\n", string(payload))
			time.Sleep(1 * time.Second)
		}
	})
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 发布消息
	publishInitMsg(client, publishTopic)

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// client.Disconnect(250)
	// fmt.Println("Disconnected")
}

func publishInitMsg(client mqtt.Client, publishTopic string) {
	payload, _ := json.Marshal(bsp.BspConfigInstance.InitConfig)
	token := client.Publish(publishTopic, 0, false, payload)
	token.Wait()
	fmt.Println("Published message:", string(payload))
	time.Sleep(1 * time.Second)
}
