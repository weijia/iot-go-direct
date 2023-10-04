package lora_shared

type EmptyArg struct {
}

type ReplyResult struct {
	Result int
}

type LoraData struct {
	Data        []byte
	ModuleIndex int // start from 0
	RSSI        float64
	SNR         float64
}
