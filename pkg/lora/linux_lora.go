//go:build linux

package lora

/*

#cgo LDFLAGS: -L${SRCDIR} -lpan3028 -lm

#include "../../pan3028-app/radio.h"
#include "../../pan3028-app/spi-3028.h"

*/
import "C"

import (
	"fmt"
	"iot_go/pkg/lora_shared"
	"os"
	"unsafe"
)

func (lora Lora) InitLora(argType lora_shared.EmptyArg, reply *lora_shared.ReplyResult) error {
	fmt.Println("Initiating Lora")
	device := C.CString(lora.DeviceName)
	defer C.free(unsafe.Pointer(device))
	C.set_device(device)
	C.rf_init()
	fmt.Println("rf_init OK\n")
	return nil
}

func (lora Lora) Exit(argType lora_shared.EmptyArg, reply *lora_shared.ReplyResult) error {
	// data := map[string]interface{}{
	// 	fmt.Sprintf("%s-exited",
	// 		lora.DeviceName): "yes",
	// }
	// bsp.GetBsp().SafeUploadTelemetry(bsp.BspConfigInstance.GatewayNodeID, data)
	reply.Result = 0
	os.Exit(0)
	reply.Result = 0
	return nil
}
