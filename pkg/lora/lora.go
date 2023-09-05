package lora

import (
	"iot_go/pkg/shared"
)

type Lora struct {
	DeviceName      string
	ModuleInst      *shared.Module
	IsInitCompleted bool
}

func NewLora(devName string) *Lora {
	return &Lora{DeviceName: devName}
}
