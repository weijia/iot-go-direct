package lora

import (
	"iot_go/pkg/shared"
)

type Lora struct {
	DeviceName string
	ModuleInst *shared.Module
}

func NewLora(devName string) *Lora {
	return &Lora{DeviceName: devName}
}
