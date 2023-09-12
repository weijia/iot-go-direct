package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMqttMsgLoop(mqttPublishCh chan interface{}, nodeMsgCh chan []byte, quit chan bool) {
	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	// MQTT 连接设置
	broker := "tcp://" + bsp.BspConfigInstance.MqttIP + ":" + fmt.Sprint(bsp.BspConfigInstance.MqttPort)
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(bsp.BspConfigInstance.GatewayNodeId)

	// 创建 MQTT 客户端
	client := mqtt.NewClient(opts)

	// 连接到服务器
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient := util.NewMqtt(&client, bsp.BspConfigInstance.GatewayNodeId)

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
		case sendingMsg := <-mqttPublishCh:
			mqttClient.SendToServer(sendingMsg)
		case receivedMsg := <-mqttCh:
			fmt.Printf("Received message: %s from topic: %s\n", receivedMsg.Payload(), receivedMsg.Topic())
			resp := HandleMqttMsg(mqttClient, receivedMsg.Payload())
			if resp != nil {
				mqttPublishCh <- resp
			}
		case nodeMsg := <-nodeMsgCh:
			HandleNodeMsg(nodeMsg, mqttPublishCh)
		case reply := <-colorUpdateRequestTimeoutCh:
			// Find the correct reply in statesForPendingColorUpdate and send reply if exists
			for index, state := range statesForPendingColorUpdate {
				fmt.Printf("data in slice: %p, data from ch: %p", &state.Reply, reply)
				if &state.Reply == reply {
					mqttPublishCh <- reply
					// Remove current element from statesForPendingColorUpdate
					statesForPendingColorUpdate = append(statesForPendingColorUpdate[:index], statesForPendingColorUpdate[index+1:]...)
					break
				}
			}
		case configRequest := <-configReqCh:
			fmt.Printf("config request: %p\n", configRequest)
			reply := finalizeConfigReq(configRequest)
			mqttPublishCh <- reply
		case <-quit:
			break
		}
	}
}
