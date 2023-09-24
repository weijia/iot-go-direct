package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/mqtt_util"
	"iot_go/pkg/node"
	"iot_go/pkg/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttPublishCh = make(chan interface{}, 10)

var mqttCh = make(chan mqtt.Message, 10)

func StartMainMsgLoop(nodeMsgCh chan []byte, quit chan bool) {
	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	easyClient := mqtt_util.MqttEasyClient{
		MqttParams:       mqtt_util.MqttParams(bsp.BspConfigInstance.MqttParams),
		ReceivingChannel: &mqttCh,
		Topic:            topic,
		MqttClientId:     bsp.BspConfigInstance.GatewayNodeId,
	}

	easyClient.SubscribeToMqttServerWithReconnect()

	mqttClient := util.NewMqtt(&easyClient.Client, bsp.BspConfigInstance.GatewayNodeId)

	for {
		select {
		case publishingMqttMsg := <-MqttPublishCh:
			mqttClient.SendToServer(publishingMqttMsg)
		case mqttMsg := <-mqttCh:
			util.IotLogInfo(fmt.Sprintf("Received message: %s from topic: %s", mqttMsg.Payload(), mqttMsg.Topic()))
			resp := HandleMqttMsg(mqttClient, mqttMsg.Payload())
			if resp != nil {
				select {
				case MqttPublishCh <- resp:
				default:
					util.IotLogErrorStr("mqtt publish channel full when sending normal mqtt reply")
				}
			}
		case nodeMsg := <-nodeMsgCh:
			HandleNodeMsg(nodeMsg, MqttPublishCh)
		case reply := <-node.ColorUpdateRequestTimeoutCh:
			// Find the correct reply in statesForPendingColorUpdate and send reply if exists
			for index, state := range node.StatesForPendingColorUpdate {
				util.IotLogInfo(fmt.Sprintf("data in slice: %p, data from ch: %p", &state.Reply, reply))
				if &state.Reply == reply {
					select {
					case MqttPublishCh <- reply:
					default:
						util.IotLogErrorStr("mqtt publish channel full when sending color update reply")
					}
					// Remove current element from statesForPendingColorUpdate
					node.StatesForPendingColorUpdate = append(
						node.StatesForPendingColorUpdate[:index], node.StatesForPendingColorUpdate[index+1:]...)
					break
				}
			}
		// case configParams := <-node.ConfigReqCh:
		// 	util.IotLogInfo(fmt.Sprintf("config request, all nodes replied or timeout, param point: %p\n", &configParams))
		// 	reply := finalizeConfigReq(configParams)
		// 	select {
		// 	case mqttPublishCh <- reply:
		// 	default:
		// 		util.IotLogErrorStr("mqtt publish channel full when sending config reply")
		// 	}
		case <-quit:
			return
		}
	}
}
