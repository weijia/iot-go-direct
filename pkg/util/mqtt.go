package util

import (
	"encoding/json"

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
	if token.Error() != nil {
		IotLogErrWithStr("Publish error", token.Error())
	} else {
		IotLog("<<<<<<<<<<<<<<<Published message: %s, %s", topic, string(payload))
	}
}

func (client Mqtt) SendToServer(data interface{}) {
	client.SendTo(client.toServerTopic, data)
}

func (client Mqtt) SendToGateway(data interface{}) {
	client.SendTo(client.toGatewayTopic, data)
}
