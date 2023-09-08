package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMqttMsgLoop(sendingCh chan interface{}, nodeMsgCh chan []byte, quit chan bool) {
	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeID + "/in"
	// MQTT 连接设置
	broker := "tcp://" + bsp.BspConfigInstance.MqttIP + ":" + fmt.Sprint(bsp.BspConfigInstance.MqttPort)
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(bsp.BspConfigInstance.GatewayNodeID)

	// 创建 MQTT 客户端
	client := mqtt.NewClient(opts)

	// 连接到服务器
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient := util.NewMqtt(&client, bsp.BspConfigInstance.GatewayNodeID)

	mqttCh := make(chan mqtt.Message)

	// 订阅主题
	token = client.Subscribe(topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
		mqttCh <- mqttMsg
	})
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case sendingMsg := <-sendingCh:
			mqttClient.SendToServer(sendingMsg)
		case receivedMsg := <-mqttCh:
			fmt.Printf("Received message: %s from topic: %s\n", receivedMsg.Payload(), receivedMsg.Topic())
			resp := HandleMsg(mqttClient, receivedMsg.Payload())
			if resp != nil {
				sendingCh <- resp
			}
		case nodeMsg := <-nodeMsgCh:
			HandleNodeMsg(nodeMsg, sendingCh)
		case <-quit:
			break
		}
	}
}
