package bsp

import (
	"encoding/json"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/shared"
	"iot_go/pkg/thingsboard"
	"iot_go/pkg/thingsboard_shared"
	"iot_go/pkg/util"

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

func InitBoard(host string) {
	module0Client = lora_client.NewLoraClient(8866, host)
	module1Client = lora_client.NewLoraClient(8867, host)
	module2Client = lora_client.NewLoraClient(8868, host)

	module0Client.InitLora(module0Param)
	module1Client.InitLora(BspConfigInstance.InitMsgContent.Module1)
	module2Client.InitLora(BspConfigInstance.InitMsgContent.Module2)
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
	// virtualBsp.thingsboardServer.CreateDevice(nodeId)
	// virtualBsp.thingsboardServer.UploadTelemetry(nodeId, data)
}

func (virtualBsp VirtualBsp) SetModule0Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	util.IotLog("Setting Module0 Params: %v", moduleParamsInJson)

	data := map[string]interface{}{
		"module0": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.SafeUploadTelemetry(BspConfigInstance.GatewayNodeId, data)
}

func (virtualBsp VirtualBsp) SetModule1Params(moduleParams shared.Module) {
	module1Client.InitLora(moduleParams)
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	util.IotLog("Setting Module1 Params: %s\n", moduleParamsInJson)

	data := map[string]interface{}{
		"module1": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.SafeUploadTelemetry(BspConfigInstance.GatewayNodeId, data)

}

func (virtualBsp VirtualBsp) SetModule2Params(moduleParams shared.Module) {
	module2Client.InitLora(moduleParams)
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	util.IotLog("Setting Module2 Params: %s\n", moduleParamsInJson)

	data := map[string]interface{}{
		"module2": fmt.Sprintf("band: %d, factor: %d, freq: %d",
			moduleParams.Band, moduleParams.Factor, moduleParams.Freq),
	}
	virtualBsp.SafeUploadTelemetry(BspConfigInstance.GatewayNodeId, data)
}

func GetBsp() *VirtualBsp {
	return &virtualBsp
}

func (virtualBsp VirtualBsp) StopAllProcess() {
	lora_client.IsQuittingSoNoRpcReconnect = true
	module0Client.Exit()
	module1Client.Exit()
	module2Client.Exit()
}

func Lora0Send(data []byte) {
	module0Client.Send(data)
}

func GetLoraClientForNode(nodeId string) *lora_client.LoraClient {
	if IsInNodeList1(nodeId) {
		return module1Client
	} else {
		if IsInNodeList2(nodeId) {
			return module2Client
		}
	}
	return nil
}

func IsInNodeList1(nodeId string) bool {
	return util.IsStrInSlice(nodeId, BspConfigInstance.BaseConfigParams.NodeList1)
}

func IsInNodeList2(nodeId string) bool {
	return util.IsStrInSlice(nodeId, BspConfigInstance.BaseConfigParams.NodeList2)
}

func GetModule0Client() *lora_client.LoraClient {
	return module0Client
}

func GetModule1Client() *lora_client.LoraClient {
	return module1Client
}

func GetModule2Client() *lora_client.LoraClient {
	return module2Client
}
