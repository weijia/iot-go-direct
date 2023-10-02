package iot_go_test

import (
	"fmt"
	"iot_go/pkg/mqtt_util"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMqttClient(t *testing.T) {
	param := mqtt_util.MqttParams{
		MqttIP:       "app.kosglass.com",
		MqttPort:     1883,
		MqttUserName: "l8juew73i2t17wavzthg",
		MqttPwd:      "i0eprmhypu3r16g3wuuc",
	}

	ch := make(chan mqtt.Message)
	topic := "device/F12309150001/in"

	easyClient := mqtt_util.MqttEasyClient{
		MqttParams:       param,
		ReceivingChannel: &ch,
		Topic:            topic,
		MqttClientId:     "F12309150001",
	}

	easyClient.ConnectAndSubscribe()

	easyClient2 := mqtt_util.MqttEasyClient{
		MqttParams:       param,
		ReceivingChannel: &ch,
		Topic:            topic,
		MqttClientId:     "F12309150001",
	}

	easyClient2.ConnectAndSubscribe()

	time.Sleep(10 * time.Second)

	easyClient2.Client.Disconnect(10)

	time.Sleep(30 * time.Second)
	fmt.Printf("Result: %d\n", 1)
}
