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
	"iot_go/pkg/shared"
	"log"
	"os"
	"unsafe"
)

func (lora Lora) InitLora(module shared.Module) int {
	fmt.Println("Initiating Lora")
	device := C.CString(lora.DeviceName)
	defer C.free(unsafe.Pointer(device))
	C.set_device(device)
	// ret = rf_init();
	isOK := int(C.rf_init())
	if isOK != 0 {
		return isOK
	}
	log.Println("rf_init OK\n")
	C.set_freq(C.int(module.Freq))
	C.set_band(C.int(module.Band))
	C.set_factor(C.int(module.Factor))
	C.rf_set_default_para()

	return 0
}

func (lora Lora) Exit() int {
	os.Exit(0)
	return 0
}
