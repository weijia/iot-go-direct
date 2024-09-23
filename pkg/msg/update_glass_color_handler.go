package msg

import (
	"encoding/hex"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/lora_module"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	// "iot_go/pkg/util"
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
		_, err := hex.DecodeString(param.Color)
		bsp.SetRequestingColor(param.NodeId, param.Color)
		if err == nil {
			if bsp.IsInNodeList1(param.NodeId) {
				colorUpdateForModule1[param.NodeId] = param
				continue
			}
			if bsp.IsInNodeList2(param.NodeId) {
				colorUpdateForModule2[param.NodeId] = param
				continue
			}
		}
		// Node id is not control of any of the lora module
		reply.Status = append(reply.Status, shared.UpdateGlassColorParams{
			NodeId: param.NodeId,
			Color:  node.SetColorForNodeAsInvalid(param.Color),
		})
	}
	var wg sync.WaitGroup
	if len(colorUpdateForModule1) > 0 {
		wg.Add(1)
		go lora_module.Module1.UpdateGlassColorForList(colorUpdateForModule1, &wg)
	}
	if len(colorUpdateForModule2) > 0 {
		wg.Add(1)
		go lora_module.Module2.UpdateGlassColorForList(colorUpdateForModule2, &wg)
	}

	wg.Wait()

	for _, param := range colorUpdateForModule1 {
		reply.Status = append(reply.Status, param)
	}
	for _, param := range colorUpdateForModule2 {
		reply.Status = append(reply.Status, param)
	}

	// Do not update color in separate go routine
	// And we do not need to update glass color here as mostly it will not complete
	// at this time
	// for _, singleNodeState := range reply.Status {
	// 	nodeState := bsp.GetOrCreateNodeState(singleNodeState.NodeId)
	// 	for i := 0; i < len(singleNodeState.Color)/2; i++ {
	// 		nodeState.NodeReportedColor[util.GetGlassAreaFromStr(singleNodeState.Color[i*2])] =
	// 			int(singleNodeState.Color[i*2+1])
	// 	}
	// 	util.IotLog("Updated state: %v, %v", nodeState, bsp.BspConfigInstance.NodeStates)
	// }

	mqttToServer <- reply
}
