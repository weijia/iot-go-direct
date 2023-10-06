//go:build linux

package lora

/*

#cgo LDFLAGS: -L${SRCDIR} -lpan3028 -lm

#include "../../pan3028-app/radio.h"
#include "../../pan3028-app/spi-3028.h"

*/
import "C"

import (
	"bytes"
	"fmt"
	"iot_go/pkg/lora_client"
	"iot_go/pkg/lora_shared"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"log"
	"os"
	"time"
	"unsafe"
)

var quit = make(chan bool)
var quitCompleted = make(chan bool)
var send = make(chan []byte)
var recv = make(chan lora_shared.LoraData, 10)
var isLoopRunning = false

func (loraDev Lora) InitLora(module shared.Module) int {
	// Do not call InitLora too frequently, otherwise, maybe isLoopRunning flag will not be set in time
	if isLoopRunning {
		quit <- true
		<-quitCompleted
	}
	util.IotLogInfo("Initiating Lora")

	device := C.CString(loraDev.DeviceName)
	defer C.free(unsafe.Pointer(device))
	C.set_device(device)
	// ret = rf_init();
	isOK := int(C.rf_init())
	if isOK != 0 {
		return isOK
	}
	log.Println("rf_init OK\n")
	freq := module.Freq * 100 * 1000
	C.set_freq(C.int(freq))
	/*
		#define BW_62_5K                        6
		#define BW_125K                         7
		#define BW_250K                         8
		#define BW_500K                         9
	*/
	bandMap := map[int]int{
		125: 7,
		250: 8,
		500: 9,
	}
	C.set_band(C.int(bandMap[module.Band]))
	C.set_factor(C.int(module.Factor))
	util.IotLogInfo(fmt.Sprintf("freq: %d, band: %d, factor: %d\n", freq, module.Band, module.Factor))

	res := int(C.rf_set_syncword(0x12)) // 0x3>0: select page, 0x12->0xf: sync word
	if res != 0 {
		fmt.Printf("-------------rf_set_syncword return non OK\n")
	}
	C.rf_set_default_para()
	// C.rf_enter_continous_rx()
	go loraDev.MsgLoop()
	return 0
}

func (loraDev Lora) CopyFromBufferIfExists() {
	// fmt.Printf("loraDev: %s, loop started: %v, loraDev: 0x%p\n", loraDev.DeviceName, loraDev.IsHandlingLoopStarted, &loraDev)
	// fmt.Printf("loraDev: %s, %d, %v", loraDev.DeviceName, loraDev.ModuleInst.Freq, loraDev.IsHandlingLoopStarted)

	buffer := bytes.NewBuffer(make([]byte, 513))

	// Call the C function to fill the buffer
	// fmt.Printf("receive to: %p\n", &buffer.Bytes()[0])
	len := int(C.receive((*C.uchar)(unsafe.Pointer(&buffer.Bytes()[0])), C.int(buffer.Cap())))

	if len <= 0 {
		// util.IotLogInfo("----------------No data received\n")
		return
	} else {
		// util.IotLogInfo(fmt.Sprintf("----------------Data received, %v\n", buffer))
		byteSlice := make([]byte, len)
		buffer.Read(byteSlice)
		args := lora_shared.LoraData{
			Data: byteSlice,
			RSSI: float64(C.RSSI),
			SNR:  float64(C.SNR),
		}

		select {
		case recv <- args:
		default:
			log.Printf("Lora.CopyFromBufferIfExists: Failed to send data to recv channel, recv channel full\n")
		}
	}
}

func (loraDev Lora) MsgLoop() {
	isLoopRunning = true
	// Create a ticker with 50 ms = 0.05 seconds
	ticker := time.NewTicker(time.Millisecond * 50)

	for {
		select {
		case <-quit:
			isLoopRunning = false
			ticker.Stop()
			quitCompleted <- true
			return

		case data := <-send:
			// Ref: https://packagewjx.github.io/2018/09/19/cgo-cstring-ram-leak/
			cBufferNeedToFree := C.CBytes(data)
			defer C.free(unsafe.Pointer(cBufferNeedToFree))
			fmt.Printf("Lora.Send: Sending len: %d\n", len(data))
			res := C.send((*C.uchar)(cBufferNeedToFree), C.int(len(data)))
			if res != 0 {
				fmt.Printf("Lora.Send: Error sending: %d\n", res)
			}
		case <-ticker.C:
			// case <-time.After(time.Second * 1):
			// util.IotLogInfo("Start event loop\n")
			C.event_handler()
			loraDev.CopyFromBufferIfExists()
			C.init_tx_or_rx()
		}
	}
	isLoopRunning = false
}

func (loraDev Lora) Exit() int {
	os.Exit(0)
	return 0
}

func (loraDev Lora) ToggleDebug() {
	C.toggle_debug()
}

func (loraDev Lora) Send(data []byte) int {
	// The following is not working as both are 0
	// util.IotLogInfo(fmt.Sprintf("len of ch %d vs cap of ch %d\n", len(send), cap(send)))

	select {
	case send <- data:
	default:
		util.IotLogInfo(fmt.Sprintf("Send buffer full, this may possiblly due to receive RPC is not started. please check the reason, %v\n", data))
		return -1
	}
	return 0
}

// func (loraDev Lora) Receive() []byte {
// 	select {
// 	case data := <-recv:
// 		return data
// 	default:
// 		return nil
// 	}
// }

func PushLoraMsgToRpc(port int, moduleIndex int, host ...string) {
	realHost := "127.0.0.1"
	if len(host) > 0 {
		realHost = host[0]
	}
	var serverRpc = lora_client.NewLoraClient(port, realHost)
	util.IotLogInfo(fmt.Sprintf("Connected to lora receive service: %s:%d", realHost, port))
	for {
		// util.IotLogInfo(fmt.Sprintf("Receiving\n"))
		select {
		case data := <-recv:
			// util.IotLogInfo(fmt.Sprintf("Pushing data: %v\n", data))
			data.ModuleIndex = moduleIndex
			serverRpc.OnReceive(data)
			// util.IotLogInfo(fmt.Sprintf("After pushing data: %v\n", data))
		}
	}
	// util.IotLogInfo(fmt.Sprintf("Quitting push msg loop\n"))
}
