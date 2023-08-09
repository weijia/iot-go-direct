package main

import (
	"encoding/json"
	"fmt"
	"iot_go/pkg/msg"
	"iot_go/pkg/util"
	"os"
	"os/signal"
	"time"

	"github.com/carlmjohnson/versioninfo"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fmt.Println("Version:", versioninfo.Version)
	fmt.Println("Revision:", versioninfo.Revision)
	config := util.FixedIotConfig{Broker: "tcp://115.159.53.168:1883", ClientId: "mqtt_golang_example"}
	deviceId := "test"
	var topic = "device/" + deviceId + "/in"
	var publishTopic = "device/" + deviceId + "/out"

	// MQTT 连接设置
	opts := mqtt.NewClientOptions().AddBroker(config.GetBroker())
	opts.SetClientID(config.GetClientId())

	// 创建 MQTT 客户端
	client := mqtt.NewClient(opts)

	// 连接到服务器
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	token = client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
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

	client.Disconnect(250)
	fmt.Println("Disconnected")
}

func publishInitMsg(client mqtt.Client, publishTopic string) {
	init := &msg.Init{
		MsgType:       "init",
		GatewayNodeID: "testing",
	}
	payload, _ := json.Marshal(init)
	token := client.Publish(publishTopic, 0, false, payload)
	token.Wait()
	fmt.Println("Published message:", payload)
	time.Sleep(1 * time.Second)
}
