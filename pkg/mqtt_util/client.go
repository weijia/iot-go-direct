package mqtt_util

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttEasyClient struct {
	MqttParams
	MqttClientId                   string
	Topic                          string
	ReceivingChannel               *chan mqtt.Message
	Client                         mqtt.Client
	IsReconnecting                 bool
	IsMqttConnectionReadyWaitGroup *sync.WaitGroup
}

func (easyClient *MqttEasyClient) createClientOptions() *mqtt.ClientOptions {

	// MQTT 连接设置
	broker := "tcp://" + easyClient.MqttIP + ":" + fmt.Sprint(easyClient.MqttPort)
	// util.IotLogInfo(fmt.Sprintf("Connecting to server %s\n", broker))
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetClientID(easyClient.MqttClientId)
	opts.Username = easyClient.MqttUserName
	opts.Password = easyClient.MqttPwd
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		util.IotLogErrWithStr("!!!!!! mqtt connection lost error", err)
		bsp.TurnOffLed("sys-led-net")
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		util.IotLog("...... mqtt reconnecting ......")
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		// 连接被建立后的回调函数
		util.IotLog("Mqtt is connected!")
		easyClient.OnConnected()
		bsp.TurnOnLed("sys-led-net")
	})
	// opts.SetConnectionLostHandler(easyClient.ReconnectMqtt)
	return opts
}

func (easyClient *MqttEasyClient) Subscribe() error {
	// 重新连接成功后重新订阅主题
	token := easyClient.Client.Subscribe(easyClient.Topic, 0, func(client mqtt.Client, mqttMsg mqtt.Message) {
		util.IotLog("Received message: %s from topic: %s", mqttMsg.Payload(), mqttMsg.Topic())
		*easyClient.ReceivingChannel <- mqttMsg
	})

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("Subscribed to topic %s\n", easyClient.Topic)
	if easyClient.IsMqttConnectionReadyWaitGroup != nil {
		easyClient.IsMqttConnectionReadyWaitGroup.Done()
	}
	return nil
}

func (easyClient *MqttEasyClient) OnConnected() {
	util.IotLog("Connected, subscribe to topic...")
	easyClient.Subscribe()
}

func (easyClient *MqttEasyClient) createClient(opts *mqtt.ClientOptions) {
	util.IotLog("Creating client for %p\n", easyClient)
	easyClient.Client = mqtt.NewClient(opts)
}

func (easyClient *MqttEasyClient) ConnectAndSubscribe() error {
	if easyClient.Client == nil {
		opts := easyClient.createClientOptions()
		easyClient.createClient(opts)
	}
	retryCnt := 0
	for {
		util.IotLog("Connecting...")
		token := easyClient.Client.Connect()

		if token.Wait() && token.Error() != nil {
			// Connect to server error
			util.IotLogErrWithStr("First connect error, need to retry", token.Error())
			// return token.Error()
			retryCnt += 1
			if retryCnt > 100 {
				log.Fatal("Can not connect to server after retry 100 times")
			}
		} else {
			easyClient.IsMqttConnectionReadyWaitGroup.Add(1)
			return nil
		}
	}
}

// func (easyClient *MqttEasyClient) ReconnectMqtt(client mqtt.Client, err error) {
// 	log.Printf("Connection lost because of: %v for %p\n", err, easyClient)
// 	if easyClient.IsReconnecting {
// 		log.Printf("Already reconnecting, for %p, client: %p should equal to %p\n", easyClient, &client, &easyClient.Client)
// 		return
// 	}
// 	easyClient.IsReconnecting = true
// 	// 在此处添加自动重连逻辑
// 	for {
// 		err := easyClient.ConnectAndSubscribe()
// 		if err != nil {
// 			log.Printf("重新连接失败, reason: %v，稍后重新尝试...\n", err)
// 			time.Sleep(15 * time.Second)
// 		} else {
// 			easyClient.IsReconnecting = false
// 			break
// 		}
// 	}
// }

// func (easyClient *MqttEasyClient) SubscribeToMqttServerWithReconnect() *mqtt.Client {
// 	for {
// 		err := easyClient.ConnectAndSubscribe()
// 		if err != nil {
// 			log.Printf("Connect and subscribe err: %v\n", err)
// 		} else {
// 			log.Printf("已连接到MQTT代理服务器并订阅主题\n")
// 			break
// 		}
// 	}
// 	return &easyClient.Client
// }
