package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MqttPublishCh = make(chan interface{}, 10)
var connectionLostHandler mqtt.ConnectionLostHandler

var mqttCh = make(chan mqtt.Message, 10)

func createClientOptions() *mqtt.ClientOptions {

	// MQTT 连接设置
	broker := "tcp://" + bsp.BspConfigInstance.MqttIP + ":" + fmt.Sprint(bsp.BspConfigInstance.MqttPort)
	util.IotLogInfo(fmt.Sprintf("Connecting to server %s\n", broker))
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(bsp.BspConfigInstance.GatewayNodeId)
	opts.Username = bsp.BspConfigInstance.MqttUserName
	opts.Password = bsp.BspConfigInstance.MqttPwd
	opts.SetConnectionLostHandler(connectionLostHandler)
	return opts
}

func ReconnectMqtt(client mqtt.Client, err error) {
	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	fmt.Println("连接丢失:", err.Error())
	// 在此处添加自动重连逻辑
	for {
		// 重新连接到MQTT代理服务器
		opts := createClientOptions()
		client := createClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			// 重新连接成功后重新订阅主题
			token = client.Subscribe(topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
				mqttCh <- mqttMsg
			})
			if token.Wait() && token.Error() != nil {
				// panic(token.Error())
				continue
			}
			fmt.Println("已重新连接到MQTT代理服务器并重新订阅主题")
			break
		}
		fmt.Println("重新连接失败，稍后重新尝试...")
		time.Sleep(5 * time.Second)
	}
}

func createClient(opts *mqtt.ClientOptions) mqtt.Client {
	client := mqtt.NewClient(opts)
	return client
}

func StartMainMsgLoop(nodeMsgCh chan []byte, quit chan bool) {
	connectionLostHandler = ReconnectMqtt
	broker := "tcp://" + bsp.BspConfigInstance.MqttIP + ":" + fmt.Sprint(bsp.BspConfigInstance.MqttPort)
	var topic = "device/" + bsp.BspConfigInstance.GatewayNodeId + "/in"
	opts := createClientOptions()
	client := createClient(opts)

	// 连接到服务器
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		//Post an device issue to thingsboard server
		//Create data string map
		data := make(map[string]interface{})
		data["fail-server"] = broker
		data["fail-username"] = opts.Username
		data["fail-password"] = opts.Password

		bsp.GetBsp().SafeUploadTelemetry(bsp.BspConfigInstance.GatewayNodeId, data)
		panic(token.Error())
	}

	// TODO: reconnect to mqtt if disconnected
	mqttClient := util.NewMqtt(&client, bsp.BspConfigInstance.GatewayNodeId)

	// 订阅主题
	token = client.Subscribe(topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
		mqttCh <- mqttMsg
	})
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

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
