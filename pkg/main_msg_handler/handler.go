package main_msg_handler

import (
	"context"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/lora_rpc"
	"iot_go/pkg/lora_shared"
	"iot_go/pkg/mqtt_util"
	"iot_go/pkg/msg"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"os/exec"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MainMsgHandler struct {
	MqttToServerCh   chan interface{}
	MqttFromServerCh chan mqtt.Message
	NodeMsgCh        chan lora_shared.LoraData
	MqttEasyClient   *mqtt_util.MqttEasyClient
	TimeoutNodeIdCh  chan string
}

func InfiniteAppLoop(ctx context.Context, loraServiceIp string) {
	for {
		h := MainMsgHandler{
			MqttToServerCh:   make(chan interface{}, 10),
			MqttFromServerCh: make(chan mqtt.Message, 10),
			NodeMsgCh:        make(chan lora_shared.LoraData, 10),
			TimeoutNodeIdCh:  make(chan string, 10),
		}
		h.TopLevelMsgLoop(ctx, loraServiceIp)
		msg.IsMsgLoopRestartNeeded = false
		msg.IsInitDone = false
		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()
		err := lora_rpc.ReceivingServer.Shutdown(ctx)
		if err != nil {
			util.IotLogErrorStr("shutting down: " + err.Error())
		} else {
			util.IotLog("shutdown processed successfully")
		}
	}
}

func (mainMsgHandler MainMsgHandler) TopLevelMsgLoop(ctx context.Context, loraServiceIp string) {
	util.IotLog("Starting main app version: %s", bsp.SwVersion)
	runLedCh := make(chan int)

	runLed := bsp.LedManager{
		DeviceName:  "sys-led-run",
		Timeout:     util.TO_SERVER_HEARTBEAT_SECONDS * 2,
		HeartbeatCh: &runLedCh,
	}

	// The application data structure can only be changed in this routine to
	// avoid concurrent data change issue
	bsp.InitConfig()

	bsp.InitBoard(loraServiceIp)
	// GetModule0Client etc will return nil before InitBoard
	lora_module.Module0.LoraClient = bsp.GetModule0Client()
	lora_module.Module1.LoraClient = bsp.GetModule1Client()
	lora_module.Module2.LoraClient = bsp.GetModule2Client()

	// Used to report node heartbeat timeout from heartbeat generator
	lora_module.TimeoutNodeIdCh = &mainMsgHandler.TimeoutNodeIdCh

	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	var wg sync.WaitGroup
	mainMsgHandler.MqttEasyClient = &mqtt_util.MqttEasyClient{
		MqttParams:                     mqtt_util.MqttParams(bsp.BspConfigInstance.MqttParams),
		ReceivingChannel:               &mainMsgHandler.MqttFromServerCh,
		Topic:                          topic,
		MqttClientId:                   bsp.BspConfigInstance.GatewayNodeId,
		IsMqttConnectionReadyWaitGroup: &wg,
	}

	mainMsgHandler.MqttEasyClient.ConnectAndSubscribe()

	mqttClient := util.NewMqtt(&mainMsgHandler.MqttEasyClient.Client,
		bsp.BspConfigInstance.GatewayNodeId)

	ctx, cancel := context.WithCancel(ctx)

	go runLed.LedMsgLoop(ctx)

	go lora_module.Module0.MsgLoop(ctx)
	go lora_module.Module1.MsgLoop(ctx)
	go lora_module.Module2.MsgLoop(ctx)
	go lora_rpc.StartLoraReceiverRpc(&mainMsgHandler.NodeMsgCh, 8869)

	// Wait for mqtt subscribe done
	wg.Wait()
	mainMsgHandler.MqttEasyClient.IsMqttConnectionReadyWaitGroup = nil // Do not need to set Wg again for reconnect

	// 发布消息
	var initMsg shared.Init
	initMsg.InitMsgContent = bsp.BspConfigInstance.InitMsgContent
	initMsg.MsgType = "init"
	mainMsgHandler.MqttToServerCh <- initMsg

	ticker := time.NewTicker(time.Second * util.TO_SERVER_HEARTBEAT_SECONDS)
	u := msg.HeartbeatStatusUpdate{
		MsgType:       "heartbeat_status_update",
		GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
	}

	for {
		select {
		case publishingMqttMsg := <-mainMsgHandler.MqttToServerCh:
			mqttClient.SendToServer(publishingMqttMsg)
			// util.IotLogInfo("Sent mqtt msg to server")
			runLedCh <- 1
			// util.IotLogInfo("Sent run led heartbeat")
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
			reply := msg.HandleMqttMsg(ctx, mainMsgHandler.MqttToServerCh, mqttMsg.Payload())
			runLedCh <- 1
			if reply != nil {
				util.SendMsgWithoutBlockingCommon(reply, mainMsgHandler.MqttToServerCh,
					fmt.Sprintf("Sending to mqtt server ch failed: %v", reply))
			}
		case nodeMsg := <-mainMsgHandler.NodeMsgCh:
			msg.HandleNodeMsg(nodeMsg, mainMsgHandler.MqttToServerCh)
		case timeoutNodeId := <-mainMsgHandler.TimeoutNodeIdCh:
			bsp.GetOrCreateNodeState(timeoutNodeId)
		// case reply := <-node.ColorUpdateRequestTimeoutCh:
		// 	// Find the correct reply in statesForPendingColorUpdate and send reply if exists
		// 	for index, state := range node.StatesForPendingColorUpdate {
		// 		util.IotLogInfo(fmt.Sprintf("data in slice: %p, data from ch: %p", &state.Reply, reply))
		// 		if &state.Reply == reply {
		// 			select {
		// 			case mainMsgHandler.MqttToServerCh <- reply:
		// 			default:
		// 				util.IotLogErrorStr("mqtt publish channel full when sending color update reply")
		// 			}
		// 			// Remove current element from statesForPendingColorUpdate
		// 			node.StatesForPendingColorUpdate = append(
		// 				node.StatesForPendingColorUpdate[:index], node.StatesForPendingColorUpdate[index+1:]...)
		// 			break
		// 		}
		// 	}
		case <-ticker.C:
			currentTimestamp := time.Now().Unix()
			l := []msg.HeartbeatStatus{}
			for _, state := range bsp.BspConfigInstance.NodeStates {
				if bsp.IsInNodeList1(state.NodeId) || bsp.IsInNodeList2(state.NodeId) {
					util.IotLog("Reporting state: %v", state)
					color := msg.GetColorStrFromSlice(state.NodeReportedColor)
					if state.LastMsgTimestamp < currentTimestamp-int64(bsp.BspConfigInstance.Heartbeat) {
						color = node.SetColorForNodeAsInvalid(color)
					}
					s := msg.HeartbeatStatus{
						NodeId:      state.NodeId,
						Color:       color,
						HardVersion: fmt.Sprintf("%x", state.HwVer),
						SoftVersion: fmt.Sprintf("%x", state.SwVer),
						RunArea:     state.RunningArea,
						RSSI:        state.RSSI,
						SNR:         state.SNR,
					}
					l = append(l, s)
				}
			}
			u.Status = l
			mqttClient.SendToServer(u)
			runLedCh <- 1
		case <-ctx.Done():
			cancel()
			return
		}
	}
}
