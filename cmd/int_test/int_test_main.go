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
	deviceId := bsp.BspConfigInstance.InitMsgContent.GatewayNodeId
	var topic = "device/" + deviceId + "/out"
	var publishTopic = "device/" + deviceId + "/in"

	// MQTT 连接设置
	broker := "tcp://" + bsp.BspConfigInstance.MqttParams.MqttIP + ":" +
		fmt.Sprint(bsp.BspConfigInstance.MqttParams.MqttPort)
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(bsp.BspConfigInstance.InitMsgContent.GatewayNodeId + "-server")
	// opts.SetClientID(bsp.BspConfigInstance.InitMsgContent.GatewayNodeId)
	// opts.Username = bsp.BspConfigInstance.MqttUserName
	// opts.Password = bsp.BspConfigInstance.MqttPwd

	// 创建 MQTT 客户端
	client := mqtt.NewClient(opts)

	// 连接到服务器
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	token = client.Subscribe(topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
		util.IotLogInfo(fmt.Sprintf("Received message: %s from topic: %s", mqttMsg.Payload(), mqttMsg.Topic()))
		// 发布消息
		var reply msg.BaseReply
		if err := json.Unmarshal(mqttMsg.Payload(), &reply); err != nil {
			util.IotLogError(err)
		}
		switch reply.MsgType {
		case "init":
			module1 := shared.Module{
				Freq:   4723,
				Band:   250,
				Factor: 10,
			}
			module2 := shared.Module{
				Freq:   4483,
				Band:   250,
				Factor: 10,
			}
			configParams := shared.ConfigParams{
				BaseConfigParams: shared.BaseConfigParams{
					// NodeList1:      []string{"FD000001", "FD000002"},
					NodeList1: []string{"01020304"},
					NodeList2: []string{"01020305"},
					// NodeList2:      []string{"01020305", "test4"},
					TouchNodeList1: []string{"01020306"},
					// TouchNodeList1: []string{"01020307", "touch2"},
					TouchNodeList2: []string{"01020307"},
					// TouchNodeList2: []string{"touch3", "touch4"},
					Custom:    "test_custom",
					Project:   "test_project",
					Heartbeat: 30,
				},
				Module1: module1,
				Module2: module2,
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
		case "config_reply":
			nodeListRequest := msg.BaseRequest{
				Method: "node_list_request",
			}
			payload, _ := json.Marshal(nodeListRequest)
			token := client.Publish(publishTopic, 0, false, payload)
			token.Wait()
			fmt.Println("Published message:", string(payload))
			time.Sleep(1 * time.Second)
		case "node_list_reply":
			nodeInfoReq := msg.GatewayNodeIdRequest{
				Method:        "node_info_request",
				GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
			}
			payload, _ := json.Marshal(nodeInfoReq)
			token := client.Publish(publishTopic, 0, false, payload)
			token.Wait()
			fmt.Println("Published message:", string(payload))
			time.Sleep(1 * time.Second)
		case "node_info_request":
			req := msg.BaseRequest{
				Method: "gateway_reboot",
			}
			payload, _ := json.Marshal(req)
			token := client.Publish(publishTopic, 0, false, payload)
			token.Wait()
			fmt.Println("Published message:", string(payload))
			time.Sleep(1 * time.Second)
		case "gateway_reboot_reply":
			params := []shared.UpdateGlassColorParams{
				{
					NodeId: "FD000001",
					Color:  "1223",
				},
				{
					NodeId: "FD000002",
					Color:  "122332",
				},
			}

			req := msg.UpdateGlassColorRequest{
				Method: "update_glass_color_request",
				Params: params,
			}
			payload, _ := json.Marshal(req)
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
