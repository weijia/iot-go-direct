package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
)

type MqttConfigRequest struct {
	Method string            `json:"method"`
	Params shared.MqttParams `json:"params"`
}

var IsMsgLoopRestartNeeded = false

func (config MqttConfigRequest) handle() interface{} {
	// fmt.Printf("%s", config.Method)
	bsp.BspConfigInstance.MqttParams.MqttIP = config.Params.MqttIP
	bsp.BspConfigInstance.MqttParams.MqttPort = config.Params.MqttPort
	bsp.BspConfigInstance.MqttParams.MqttUserName = config.Params.MqttUserName
	bsp.BspConfigInstance.MqttParams.MqttPwd = config.Params.MqttPwd
	bsp.BspConfigInstance.CommitChanges()
	IsRebootNeeded = true
	var mqttConfigReply BaseReply
	mqttConfigReply.MsgType = "mqtt_config_reply"
	return mqttConfigReply
}
