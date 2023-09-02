//go:build windows

package lora

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_shared"
	"os"
)

func (lora Lora) InitLora(argType lora_shared.EmptyArg, reply *lora_shared.ReplyResult) error {
	fmt.Println("Initiating Lora")
	// C.rf_init()
	data := map[string]interface{}{
		fmt.Sprintf("%s-inited",
			lora.DeviceName): "yes",
	}
	bsp.GetBsp().SafeUploadTelemetry(bsp.BspConfigInstance.GatewayNodeID, data)
	reply.Result = 0
	return nil
}

func (lora Lora) Exit(argType lora_shared.EmptyArg, reply *lora_shared.ReplyResult) error {
	data := map[string]interface{}{
		fmt.Sprintf("%s-exited",
			lora.DeviceName): "yes",
	}
	bsp.GetBsp().SafeUploadTelemetry(bsp.BspConfigInstance.GatewayNodeID, data)
	reply.Result = 0
	os.Exit(0)
	reply.Result = 0
	return nil
}
