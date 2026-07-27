package main_msg_handler

import (
	"context"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
	"iot_go/pkg/mqtt_util"
	"iot_go/pkg/msg"
	"iot_go/pkg/node"
	"iot_go/pkg/serial"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"os/exec"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	NODE_STATE_CMD_COMPLETED = 2
)

type MainMsgHandler struct {
	MqttToServerCh   chan interface{}
	MqttFromServerCh chan mqtt.Message
	MqttEasyClient   *mqtt_util.MqttEasyClient
	TimeoutNodeIdCh  chan string
}

// InfiniteAppLoop 持续运行主消息循环。port 为已打开的串口(真实或 mock)；
// skipMqtt 为 true 时(例如离线 mock 测试)不连接 MQTT broker，仅跑控制板通信。
func InfiniteAppLoop(ctx context.Context, port *serial.Port, skipMqtt bool) {
	for {
		h := MainMsgHandler{
			MqttToServerCh:   make(chan interface{}, 10),
			MqttFromServerCh: make(chan mqtt.Message, 10),
			TimeoutNodeIdCh:  make(chan string, 10),
		}
		h.TopLevelMsgLoop(ctx, port, skipMqtt)
		msg.IsMsgLoopRestartNeeded = false
		msg.IsInitDone = false
		util.IotLog("shutdown processed successfully")
	}
}

func (mainMsgHandler MainMsgHandler) TopLevelMsgLoop(ctx context.Context, port *serial.Port, skipMqtt bool) {
	util.IotDebugPrintf("Starting main app version: %s", bsp.SwVersion)
	runLedCh := make(chan int)

	runLed := bsp.LedManager{
		DeviceName:  "sys-led-run",
		Timeout:     util.TO_SERVER_HEARTBEAT_SECONDS * 2,
		HeartbeatCh: &runLedCh,
	}

	// 应用数据只在主循环(及其派生的同步调用)里修改，避免并发数据竞争
	bsp.InitConfig()

	// 单串口控制板替代原 3 个 LoRa 模块 + RPC 服务
	controller.Init(port)
	go controller.Ctrl.MsgLoop(ctx)
	// 读协程收到板子帧后，转发给 MsgLoop 里正在等待的发送者
	go func() {
		for f := range port.RecvCh() {
			controller.Ctrl.RecvCh() <- f
		}
	}()

	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	var wg sync.WaitGroup
	var mqttClient *util.Mqtt

	if !skipMqtt {
		mainMsgHandler.MqttEasyClient = &mqtt_util.MqttEasyClient{
			MqttParams:                     mqtt_util.MqttParams(bsp.BspConfigInstance.MqttParams),
			ReceivingChannel:               &mainMsgHandler.MqttFromServerCh,
			Topic:                          topic,
			MqttClientId:                   bsp.BspConfigInstance.GatewayNodeId,
			IsMqttConnectionReadyWaitGroup: &wg,
		}
		mainMsgHandler.MqttEasyClient.ConnectAndSubscribe()
		mqttClient = util.NewMqtt(&mainMsgHandler.MqttEasyClient.Client,
			bsp.BspConfigInstance.GatewayNodeId)
	}

	ctx, cancel := context.WithCancel(ctx)

	go runLed.LedMsgLoop(ctx)

	// Wait for mqtt subscribe done
	if !skipMqtt {
		wg.Wait()
		mainMsgHandler.MqttEasyClient.IsMqttConnectionReadyWaitGroup = nil // Do not need to set Wg again for reconnect
	}

	// 发布消息
	var initMsg shared.Init
	initMsg.InitMsgContent = bsp.BspConfigInstance.InitMsgContent
	initMsg.MsgType = "init"
	mainMsgHandler.MqttToServerCh <- initMsg

	// Use 10 seconds as period
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	u := msg.HeartbeatStatusUpdate{
		MsgType:       "heartbeat_status_update",
		GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
	}

	for {
		select {
		case publishingMqttMsg := <-mainMsgHandler.MqttToServerCh:
			util.IotDebug("MainLoop: Sending mqtt msg to server")
			if mqttClient != nil {
				mqttClient.SendToServer(publishingMqttMsg)
			} else {
				util.IotLogInfo(fmt.Sprintf("MainLoop: skip publishing (mock/no mqtt): %v", publishingMqttMsg))
			}
			runLedCh <- 1
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
				if mainMsgHandler.MqttEasyClient != nil {
					mainMsgHandler.MqttEasyClient.Client.Disconnect(1000)
				}
				cancel()
			}
		case mqttMsg := <-mainMsgHandler.MqttFromServerCh:
			util.IotDebug("MainLoop: received mqtt msg from server")
			reply := msg.HandleMqttMsg(ctx, mainMsgHandler.MqttToServerCh, mqttMsg.Payload())
			runLedCh <- 1
			if reply != nil {
				util.SendMsgWithoutBlockingCommon(reply, mainMsgHandler.MqttToServerCh,
					fmt.Sprintf("Sending to mqtt server ch failed: %v", reply))
			}
		case timeoutNodeId := <-mainMsgHandler.TimeoutNodeIdCh:
			util.IotDebug("MainLoop: Handle node msg timeout")
			nodeState := bsp.GetOrCreateNodeState(timeoutNodeId)
			nodeState.IsOffline = true

		case <-ticker.C:
			bsp.RemainingPeriod -= 1
			if bsp.RemainingPeriod <= 0 {
				bsp.RemainingPeriod = bsp.PeriodNumberForReportingToServer
			} else {
				continue
			}
			util.IotDebug("MainLoop: prepare heartbeat to server")
			l := []msg.HeartbeatStatus{}
			resendCmd := false
			for _, state := range bsp.BspConfigInstance.NodeStates {
				if bsp.IsInNodeList1(state.NodeId) || bsp.IsInNodeList2(state.NodeId) {
					util.IotDebugPrintf("Reporting state: %v", state)
					color := msg.GetColorStrFromSlice(state.NodeReportedColor)
					requestingColor := msg.GetColorStrFromSlice(state.NodeRequestingColor)
					if state.IsOffline {
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
						IsOffline:   state.IsOffline,
						CompletionStatus: state.CompletionStatus,
						NodeRequestingColor: requestingColor,
					}
					l = append(l, s)
					if !bsp.IsEqual(state.NodeReportedColor, state.NodeRequestingColor) {
						resendCmd = true
					}
				}
			}
			if resendCmd {
				msg.KeptGroupUpdateGlassColorRequest.Replay()
			}
			u.Status = l
			if mqttClient != nil {
				mqttClient.SendToServer(u)
			}
			runLedCh <- 1
		// Handle ctx cancel
		case <-ctx.Done():
			cancel()
			return
		}
	}
}
