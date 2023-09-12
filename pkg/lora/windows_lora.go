//go:build windows

package lora

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"os"
)

func (lora Lora) InitLora(module shared.Module) int {
	fmt.Println("Initiating Lora")
	// C.rf_init()
	data := map[string]interface{}{
		fmt.Sprintf("%s-inited",
			lora.DeviceName): "yes",
	}
	bsp.GetBsp().SafeUploadTelemetry(bsp.BspConfigInstance.GatewayNodeId, data)
	return 0
}

func (lora Lora) Exit() int {
	data := map[string]interface{}{
		fmt.Sprintf("%s-exited",
			lora.DeviceName): "yes",
	}
	bsp.GetBsp().SafeUploadTelemetry(bsp.BspConfigInstance.GatewayNodeId, data)
	os.Exit(0)
	return 0
}

func (lora Lora) Send(data []byte) int {
	return 0
}

func (lora Lora) Receive() []byte {
	return make([]byte, 0)
}

func (lora Lora) ToggleDebug() {
}

func PushLoraMsgToRpc(pushPort int, pushHost ...string) {
	util.IotLogInfo("Started dummy windows push msg RPC\n")
}
