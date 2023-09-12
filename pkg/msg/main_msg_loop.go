package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
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
		case publishingMqttMsg := <-mqttPublishCh:
			mqttClient.SendToServer(publishingMqttMsg)
		case mqttMsg := <-mqttCh:
			util.IotLogInfo(fmt.Sprintf("Received message: %s from topic: %s\n", mqttMsg.Payload(), mqttMsg.Topic()))
			resp := HandleMqttMsg(mqttClient, mqttMsg.Payload())
			if resp != nil {
				select {
				case mqttPublishCh <- resp:
				default:
					util.IotLogErrorStr("mqtt publish channel full when sending normal mqtt reply")
				}
			}
		case nodeMsg := <-nodeMsgCh:
			HandleNodeMsg(nodeMsg, mqttPublishCh)
		case reply := <-colorUpdateRequestTimeoutCh:
			// Find the correct reply in statesForPendingColorUpdate and send reply if exists
			for index, state := range statesForPendingColorUpdate {
				util.IotLogInfo(fmt.Sprintf("data in slice: %p, data from ch: %p", &state.Reply, reply))
				if &state.Reply == reply {
					select {
					case mqttPublishCh <- reply:
					default:
						util.IotLogErrorStr("mqtt publish channel full when sending color update reply")
					}
					// Remove current element from statesForPendingColorUpdate
					statesForPendingColorUpdate = append(statesForPendingColorUpdate[:index], statesForPendingColorUpdate[index+1:]...)
					break
				}
			}
		case configParams := <-node.ConfigReqCh:
			util.IotLogInfo(fmt.Sprintf("config request: %p\n", configParams))
			reply := finalizeConfigReq(configParams)
			select {
			case mqttPublishCh <- reply:
			default:
				util.IotLogErrorStr("mqtt publish channel full when sending config reply")
			}
		case <-quit:
			break
		}
	}
}
