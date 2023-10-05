package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"sync"
)

type UpdateGlassColorRequest struct {
	Method string                          `json:"method"`
	Params []shared.UpdateGlassColorParams `json:"params"`
}

func SetSingleGlassColor(client *lora_client.LoraClient, nodeId string, color string) {
	c := node.GetUpdateGlassColorMsg(
		nodeId, color)
	client.Send(c)
}

func (request UpdateGlassColorRequest) handle(mqttToServer chan interface{}) {
	var reply shared.UpdateGlassColorReply

	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.MsgType = "update_glass_color_reply"

	colorUpdateForModule1 := make(map[string]shared.UpdateGlassColorParams)
	colorUpdateForModule2 := make(map[string]shared.UpdateGlassColorParams)
	reply.Status = make([]shared.UpdateGlassColorParams, 0, len(request.Params))

	for _, param := range request.Params {
		if bsp.IsInNodeList1(param.NodeId) {
			colorUpdateForModule1[param.NodeId] = param
		} else if bsp.IsInNodeList2(param.NodeId) {
			colorUpdateForModule2[param.NodeId] = param
		} else {
			reply.Status = append(reply.Status, shared.UpdateGlassColorParams{
				NodeId: param.NodeId,
				Color: node.SetColorForNodeAsInvalid(param.Color),
			})
		}
	}
	var wg sync.WaitGroup
	if len(colorUpdateForModule1) > 0 {
		wg.Add(1)
		go lora_module.Module1.UpdateGlassColorForList(colorUpdateForModule1, &wg)
	}
	if len(colorUpdateForModule2) > 0 {
		wg.Add(1)
		go lora_module.Module1.UpdateGlassColorForList(colorUpdateForModule2, &wg)
	}
	
	wg.Wait()

	for _, param := range colorUpdateForModule1 {
		reply.Status = append(reply.Status, param)
	}
	for _, param := range colorUpdateForModule2 {
		reply.Status = append(reply.Status, param)
	}

	mqttToServer <- reply
}
