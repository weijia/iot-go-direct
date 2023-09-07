package bsp

import (
	"encoding/json"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/node"
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

var module0Client *lora_client.LoraClient
var module1Client *lora_client.LoraClient
var module2Client *lora_client.LoraClient

func InitBoard() {
	module0Client = lora_client.NewLoraClient(8866)
	module1Client = lora_client.NewLoraClient(8867)
	module2Client = lora_client.NewLoraClient(8868)

	module0Client.InitLora(module0Param)
	module1Client.InitLora(module1Param)
	module2Client.InitLora(module2Param)
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

func (virtualBsp VirtualBsp) SafeUploadTelemetry(nodeId string, data interface{}) {
	virtualBsp.thingsboardServer.CreateDevice(nodeId)
	virtualBsp.thingsboardServer.UploadTelemetry(nodeId, data)
}

func (virtualBsp VirtualBsp) SetModule0Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Println("Setting Module0 Params", moduleParamsInJson)

	data := map[string]interface{}{
		"module0": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.SafeUploadTelemetry(BspConfigInstance.GatewayNodeID, data)
}

func (virtualBsp VirtualBsp) SetModule1Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Printf("Setting Module1 Params: %s\n", moduleParamsInJson)

	data := map[string]interface{}{
		"module1": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.SafeUploadTelemetry(BspConfigInstance.GatewayNodeID, data)
}

func (virtualBsp VirtualBsp) SetModule2Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Printf("Setting Module2 Params: %s\n", moduleParamsInJson)

	data := map[string]interface{}{
		"module2": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.SafeUploadTelemetry(BspConfigInstance.GatewayNodeID, data)
}

func IsSliceContainsStr(a []string, b string) bool {
	for _, c := range a {
		if c == b {
			return true
		}
	}
	return false
}

func (virtualBsp VirtualBsp) SetSingleGlassColor(nodeId string, color string) {
	fmt.Println("Setting glass color:", nodeId, color)

	if IsSliceContainsStr(BspConfigInstance.BaseConfigParams.NodeList1, nodeId) {
		module1Client.Send(node.GetUpdateGlassColorMsg(nodeId, color))
	} else {
		if IsSliceContainsStr(BspConfigInstance.BaseConfigParams.NodeList2, nodeId) {
			module1Client.Send(node.GetUpdateGlassColorMsg(nodeId, color))
		} else {
			util.IotLogErrorStr("Node does not exists\n")
		}
	}

	data := map[string]interface{}{
		"color": color,
	}
	virtualBsp.SafeUploadTelemetry(nodeId, data)
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

func (virtualBsp VirtualBsp) StopAllProcess() {
	module0Client.Exit()
	module1Client.Exit()
	module2Client.Exit()
}

func Lora0Send(data []byte) {
	module0Client.Send(data)
}
