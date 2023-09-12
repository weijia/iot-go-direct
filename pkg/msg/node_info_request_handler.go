package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/util"
)

type NodeInfoRequest struct {
	Method string        `json:"method"`
	Params GatewayParams `json:"params"`
}
type GatewayParams struct {
	GatewayNodeId string `json:"gateway_node_id"`
}

func (nodeInfoRequest NodeInfoRequest) handle(mqttClient *util.Mqtt) interface{} {
	var reply NodeInfoReply
	reply.MsgType = "node_info_reply"
	reply.NodeInfoContent = bsp.BspConfigInstance.InitMsgContent.NodeInfoContent
	reply.MqttParams = bsp.BspConfigInstance.MqttParams
	return reply
}
