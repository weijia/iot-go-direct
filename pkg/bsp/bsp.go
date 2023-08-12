package bsp

import (
	"encoding/json"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"time"

	"fmt"
)

type Bsp interface {
	SetModule0Params(module shared.Module)
	SetModule1Params(module shared.Module)
	SetModule2Params(module shared.Module)
}

func InitBoard() {

}

type VirtualBsp struct {
}

func (virtualBsp VirtualBsp) SetModule0Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Println("Setting Module0 Params", moduleParamsInJson)
}

func (virtualBsp VirtualBsp) SetModule1Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Printf("Setting Module1 Params: %s\n", moduleParamsInJson)
}

func (virtualBsp VirtualBsp) SetModule2Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Printf("Setting Module2 Params: %s\n", moduleParamsInJson)
}

func (virtualBsp VirtualBsp) SetSingleGlassColor(nodeId string, color string) {
	fmt.Println("Setting glass color:", nodeId, color)
}

func (virtualBsp VirtualBsp) SetGlassColorsBlocking(
	mqttClient *util.Mqtt, updateGlassColorParams []shared.UpdateGlassColorParams) {
	// using for loop
	for _, value := range updateGlassColorParams {
		virtualBsp.SetSingleGlassColor(value.NodeID, value.Color)
	}
	time.Sleep(90 * time.Second)
	var reply shared.UpdateGlassColorReply

	reply.GatewayNodeID = BspConfigInstance.GatewayNodeID
	reply.MsgType = "update_glass_color_reply"
	reply.Status = updateGlassColorParams
	mqttClient.SendToServer(reply)
}

func GetBsp() VirtualBsp {
	var virtualBsp VirtualBsp
	return virtualBsp
}
