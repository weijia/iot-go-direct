//go:build linux

package lora

/*

#cgo LDFLAGS: -L${SRCDIR} -lpan3028 -lm

#include "../../pan3028-app/radio.h"

*/
import "C"

import (
	"fmt"
	"iot_go/pkg/lora_shared"
)

func (lora Lora) InitLora(argType lora_shared.EmptyArg, reply *lora_shared.ReplyResult) error {
	fmt.Println("Initiating Lora")
	C.rf_init()
}
