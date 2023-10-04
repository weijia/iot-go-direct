package main_msg_handler

import (
	"context"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/lora_shared"
	"iot_go/pkg/mqtt_util"
	"iot_go/pkg/msg"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MainMsgHandler struct {
	MqttToServerCh   chan interface{}
	MqttFromServerCh chan mqtt.Message
	NodeMsgCh chan lora_shared.LoraData
	MqttEasyClient *mqtt_util.MqttEasyClient
}


func InfiniteAppLoop(ctx context.Context, loraServiceIp string) {
	for {
		h := MainMsgHandler{
			MqttToServerCh: make(chan interface{}, 10),
			MqttFromServerCh: make(chan mqtt.Message, 10),
			NodeMsgCh: make(chan lora_shared.LoraData, 10),
		}
		h.TopLevelMsgLoop(ctx, loraServiceIp)
		msg.IsMsgLoopRestartNeeded = false
		msg.IsInitDone = false
	}
}


func (mainMsgHandler MainMsgHandler) TopLevelMsgLoop(ctx context.Context, loraServiceIp string) {
	bsp.InitConfig()
	// fmt.Println(viper.GetString("msg_type"))

	bsp.InitBoard(loraServiceIp)
	util.IotLog("Starting main app version: %s", bsp.SwVersion)

	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	mainMsgHandler.MqttEasyClient = &mqtt_util.MqttEasyClient{
		MqttParams:       mqtt_util.MqttParams(bsp.BspConfigInstance.MqttParams),
		ReceivingChannel: &mainMsgHandler.MqttFromServerCh,
		Topic:            topic,
		MqttClientId:     bsp.BspConfigInstance.GatewayNodeId,
	}

	mainMsgHandler.MqttEasyClient.ConnectAndSubscribe()

	mqttClient := util.NewMqtt(&mainMsgHandler.MqttEasyClient.Client, 
		bsp.BspConfigInstance.GatewayNodeId)

	ctx, cancel := context.WithCancel(ctx)
	
	go lora_module.Module0.MsgLoop(ctx)
	go lora_module.Module1.MsgLoop(ctx)
	go lora_module.Module2.MsgLoop(ctx)

	go msg.SendHeartbeatToServer(ctx, mainMsgHandler.MqttToServerCh)

	// 发布消息
	var initMsg shared.Init
	initMsg.InitMsgContent = bsp.BspConfigInstance.InitMsgContent
	initMsg.MsgType = "init"
	mainMsgHandler.MqttToServerCh <- initMsg

	for {
		select {
		case publishingMqttMsg := <-mainMsgHandler.MqttToServerCh:
			mqttClient.SendToServer(publishingMqttMsg)
			// We put the restart here instead of after handling mqtt msg
			// so we will only restart after msg published to server
			if msg.IsRebootNeeded {
				cmd := exec.Command("reboot")
				err := cmd.Run()
				if err != nil {
					util.IotLogErrorStr("Execute reboot failed")
				}
			}
			if msg.IsMsgLoopRestartNeeded {
				mainMsgHandler.MqttEasyClient.Client.Disconnect(1000)
				cancel()
			}
		case mqttMsg := <-mainMsgHandler.MqttFromServerCh:
			reply := msg.HandleMqttMsg(ctx, mqttClient, mqttMsg.Payload())
			if reply != nil {
				util.SendMsgWithoutBlockingCommon(reply, mainMsgHandler.MqttToServerCh,
						fmt.Sprintf("Sending to mqtt server failed: %v", reply))
			}
		case nodeMsg := <-mainMsgHandler.NodeMsgCh:
			msg.HandleNodeMsg(nodeMsg, mainMsgHandler.MqttToServerCh)
		case reply := <-node.ColorUpdateRequestTimeoutCh:
			// Find the correct reply in statesForPendingColorUpdate and send reply if exists
			for index, state := range node.StatesForPendingColorUpdate {
				util.IotLogInfo(fmt.Sprintf("data in slice: %p, data from ch: %p", &state.Reply, reply))
				if &state.Reply == reply {
					select {
					case mainMsgHandler.MqttToServerCh <- reply:
					default:
						util.IotLogErrorStr("mqtt publish channel full when sending color update reply")
					}
					// Remove current element from statesForPendingColorUpdate
					node.StatesForPendingColorUpdate = append(
						node.StatesForPendingColorUpdate[:index], node.StatesForPendingColorUpdate[index+1:]...)
					break
				}
			}

		case <-ctx.Done():
			cancel()
			return
		}
	}
}