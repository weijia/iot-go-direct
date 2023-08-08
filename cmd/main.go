package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	var mqttConfig = util.FixedIotConfig{Broker: "tcp://115.159.53.168:1883", ClientId: "mqtt_golang_example"}
	var topic = "test-go"

	// MQTT 连接设置
	opts := mqtt.NewClientOptions().AddBroker(mqttConfig.GetBroker())
	opts.SetClientID(mqttConfig.GetClientId())

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
	for i := 0; i < 5; i++ {
		payload := "Test message " + strconv.Itoa(i)
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		fmt.Println("Published message:", payload)
		time.Sleep(1 * time.Second)
	}

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	client.Disconnect(250)
	fmt.Println("Disconnected")
}
