//go:build windows

package lora

import (
	"fmt"
	"os"
)

func (lora Lora) InitLora(argType EmptyArg, reply *ReplyResult) error {
	fmt.Println("Initiating Lora")
	// C.rf_init()
	reply.Result = 0
	return nil
}

func (lora Lora) Exit(argType EmptyArg, reply *ReplyResult) error {
	os.Exit(0)
	reply.Result = 0
	return nil
}
