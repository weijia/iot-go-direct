package msg

import (
	"iot_go/pkg/shared"
)

type MqttConfigRequest struct {
	Method string            `json:"method"`
	Params shared.MqttParams `json:"params"`
}
