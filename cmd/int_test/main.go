package main

import (
	"encoding/json"
	"fmt"
	"time"

	"os"
	"os/signal"

	"iot_go/pkg/bsp"
	"iot_go/pkg/msg"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	bsp.InitConfig()
	deviceId := bsp.BspConfigInstance.InitConfig.GatewayNodeID
	var topic = "device/" + deviceId + "/out"
	var publishTopic = "device/" + deviceId + "/in"

	// MQTT 连接设置
	opts := mqtt.NewClientOptions().AddBroker(bsp.BspConfigInstance.Broker)
	opts.SetClientID(bsp.BspConfigInstance.InitConfig.GatewayNodeID + "-server")

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
		// 发布消息
		var reply msg.BaseReply
		if err := json.Unmarshal(mqttMsg.Payload(), &reply); err != nil {
			util.IotLogFatal(err)
		}
		switch reply.MsgType {
		case "init":
			module1 := shared.Module{
				Freq:   447,
				Band:   250,
				Factor: 10,
			}
			module2 := shared.Module{
				Freq:   448,
				Band:   250,
				Factor: 10,
			}
			configParams := shared.ConfigParams{
				NodeList1:      []string{"test1", "test2"},
				NodeList2:      []string{"test3", "test4"},
				TouchNodeList1: []string{"touch1", "touch2"},
				TouchNodeList2: []string{"touch3", "touch4"},
				Custom:         "test_custom",
				Project:        "test_project",
				HeartBeat:      30,
				Module1:        module1,
				Module2:        module2,
			}

			configRequest := msg.ConfigRequest{
				Method: "config",
				Params: configParams,
			}
			payload, _ := json.Marshal(configRequest)
			token := client.Publish(publishTopic, 0, false, payload)
			token.Wait()
			fmt.Println("Published message:", string(payload))
			time.Sleep(1 * time.Second)
		}

	})

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// client.Disconnect(250)
	// fmt.Println("Disconnected")
}
