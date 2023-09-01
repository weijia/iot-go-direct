package lora

import (
	"fmt"
	"iot_go/pkg/shared"
	"log"
	"net/http"
	"net/rpc"
)

type Lora struct {
	DeviceName string
	ModuleInst *shared.Module
}

func NewLora(devName string) *Lora {
	return &Lora{DeviceName: devName}
}

type EmptyArg struct {
}

type ReplyResult struct {
	Result int
}

func StartLoraService(devName string, port int) {
	fmt.Printf("Starting lora service on dev: %s, port: %d", devName, port)
	rolaDev := NewLora(devName)
	rpc.Register(rolaDev)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("serve error:", err)
	}
}
