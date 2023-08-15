package bsp

import (
	"encoding/json"
	"iot_go/pkg/shared"
	"iot_go/pkg/thingsboard"
	"iot_go/pkg/thingsboard_shared"
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
	thingsboardServer thingsboard.ThingsboardServer
}

var virtualBsp = VirtualBsp{
	thingsboardServer: thingsboard.ThingsboardServer{
		DeviceProfile: thingsboard_shared.DeviceProfile{
			ThingsboardServerInfo: thingsboard_shared.ThingsboardServerInfo{
				Server: "http://120.55.92.168",
				Port:   8080,
			},
			ProvisionKey:    "0hsh1hpc605g4kwyal46",
			ProvisionSecret: "68rsgqafhw0anhcwnccr",
		},
	},
}

func (virtualBsp VirtualBsp) SetModule0Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	virtualBsp.thingsboardServer.CreateDevice(BspConfigInstance.GatewayNodeID)
	fmt.Println("Setting Module0 Params", moduleParamsInJson)
	data := map[string]interface{}{
		"module0": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.thingsboardServer.UploadTelemetry(BspConfigInstance.GatewayNodeID, data)
}

func (virtualBsp VirtualBsp) SetModule1Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Printf("Setting Module1 Params: %s\n", moduleParamsInJson)
	virtualBsp.thingsboardServer.CreateDevice(BspConfigInstance.GatewayNodeID)
	data := map[string]interface{}{
		"module1": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.thingsboardServer.UploadTelemetry(BspConfigInstance.GatewayNodeID, data)
}

func (virtualBsp VirtualBsp) SetModule2Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Printf("Setting Module2 Params: %s\n", moduleParamsInJson)
	virtualBsp.thingsboardServer.CreateDevice(BspConfigInstance.GatewayNodeID)
	data := map[string]interface{}{
		"module2": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.thingsboardServer.UploadTelemetry(BspConfigInstance.GatewayNodeID, data)
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

func GetBsp() *VirtualBsp {
	return &virtualBsp
}
