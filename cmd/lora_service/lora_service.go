package main

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora"
)

func main() {
	bsp.InitConfig()
	lora.StartLoraService("spidev0.0", 8866)
}
