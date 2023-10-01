package util

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt struct {
	Client         *mqtt.Client
	DeviceId       string
	toServerTopic  string
	toGatewayTopic string
}

func NewMqtt(client *mqtt.Client, deviceId string) *Mqtt {
	return &Mqtt{
		Client:         client,
		DeviceId:       deviceId,
		toServerTopic:  "device/" + deviceId + "/out",
		toGatewayTopic: "device/" + deviceId + "/in",
	}
}

func (client Mqtt) SendTo(topic string, data interface{}) {
	payload, _ := json.Marshal(data)
	token := (*client.Client).Publish(topic, 0, false, payload)
	token.Wait()
	IotLog("Published message:", topic, string(payload))
	time.Sleep(1 * time.Second)
}

func (client Mqtt) SendToServer(data interface{}) {
	client.SendTo(client.toServerTopic, data)
}

func (client Mqtt) SendToGateway(data interface{}) {
	client.SendTo(client.toGatewayTopic, data)
}
