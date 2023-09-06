package lora

// The lora object will be copied to RPC procedure instead of sending the original object
type Lora struct {
	DeviceName string
	// ModuleInst            *shared.Module
	// IsInitCompleted       bool
	IsHandlingLoopStarted bool
}

func NewLora(devName string) *Lora {
	return &Lora{DeviceName: devName, IsHandlingLoopStarted: false}
}
