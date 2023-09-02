package main

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora"
)

func main() {
	bsp.InitConfig()
	lora.StartLoraService("/dev/spidev2.0", 8866)
}
