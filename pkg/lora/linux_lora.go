//go:build linux

package lora

/*

#cgo LDFLAGS: -L${SRCDIR} -lpan3028 -lm

#include "../../pan3028-app/radio.h"

*/
import "C"

import (
	"fmt"
)

func (lora Lora) InitLora() {
	fmt.Println("Initiating Lora")
	C.rf_init()
}
