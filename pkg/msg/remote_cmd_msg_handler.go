package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/remote_cmd"
)

type RemoteCmdReq struct {
	Method string            `json:"method"`
	Params RemoteCmdParams   `json:"params"`
}
type RemoteCmdParams struct {
	Cmd string `json:"cmd"`
}

type RemoteCmdReply struct {
	MsgType       string `json:"msg_type"`
	GatewayNodeId string `json:"gateway_node_id"`
	remote_cmd.RemoteCmdResult
}


func (req RemoteCmdReq) handle(mqttToServer chan interface{}) {
	var reply RemoteCmdReply
	res := remote_cmd.ExecuteCmd(req.Params.Cmd)
	reply.RemoteCmdResult = res
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.MsgType = "remote_cmd_result"
	mqttToServer <- reply
}
