package msg

import (
	"context"
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"time"
)

type UpdateGlassColorRequest struct {
	Method string                          `json:"method"`
	Params []shared.UpdateGlassColorParams `json:"params"`
}

func SetSingleGlassColor(client *lora_client.LoraClient, nodeId string, color string) {
	c := node.GetUpdateGlassColorMsg(
		util.DecodeId(bsp.BspConfigInstance.GatewayNodeId),
		util.DecodeId(nodeId), color)
	client.Send(c)
}

func (request UpdateGlassColorRequest) handle() {
	var reply shared.UpdateGlassColorReply

	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.MsgType = "update_glass_color_reply"

	// Generate a map between node Id and lora client
	// loraClientMap := make(map[string]*lora_client.LoraClient)
	updateGlassColorParams := request.Params

	ctx, cancel := context.WithCancel(context.Background())

	state := node.ColorUpdateState{
		Reply:      reply,
		Timer:      time.NewTimer(120 * time.Second),
		CancelFunc: &cancel,
	}

	for _, param := range updateGlassColorParams {
		client := bsp.GetLoraClientForNode(param.NodeId)
		if client == nil {
			util.IotLogErrorStr(fmt.Sprintf("Node: %s does not exists\n", param.NodeId))
			param.Color = node.SetColorForNodeAsInvalid(param.Color)
			reply.Status = append(reply.Status, param)
		} else {
			reply.Status = append(reply.Status, param)
			SetSingleGlassColor(client, param.NodeId, param.Color)
			state.PendingNodes = append(state.PendingNodes, param.NodeId)
		}
		bsp.GetBsp().SafeUploadTelemetry(param.NodeId+"-requesting", param.Color)
	}

	node.StatesForPendingColorUpdate = append(node.StatesForPendingColorUpdate, state)
	// Send reply pointer to requestTimeoutCh after 120 seconds
	// TODO: Maybe we need to retry before actual 120 second timeout
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-state.Timer.C:
			node.ColorUpdateRequestTimeoutCh <- &state.Reply
		}
	}()
}
