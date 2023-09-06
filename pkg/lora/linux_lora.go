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
	"iot_go/pkg/shared"
	"log"
	"os"
	"unsafe"
)

func (loraDev Lora) SendReceiveLoop() {
	fmt.Println("calling send_receive_loop\n")
	C.send_receive_loop()
	fmt.Println("after calling send_receive_loop\n")
}

func (loraDev Lora) ReceiveLoop() {
	fmt.Println("calling receive_loop\n")
	C.receive_loop()
	fmt.Println("after calling receive_loop\n")
}

func (loraDev Lora) InitLora(module shared.Module) int {
	fmt.Println("Initiating Lora")
	device := C.CString(loraDev.DeviceName)
	defer C.free(unsafe.Pointer(device))
	C.set_device(device)
	// ret = rf_init();
	isOK := int(C.rf_init())
	if isOK != 0 {
		return isOK
	}
	log.Println("rf_init OK\n")
	C.set_freq(C.int(module.Freq) * 1000 * 1000)
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
	fmt.Printf("freq: %d, band: %d, factor: %d\n", module.Freq, module.Band, module.Factor)

	return 0
}

func (loraDev Lora) Exit() int {
	os.Exit(0)
	return 0
}

func (loraDev Lora) ToggleDebug() {
	C.toggle_debug()
}

func (loraDev Lora) Send(data []byte) int {
	if C.is_loop_started() == 0 {
		fmt.Printf("Start handling loop\n")
		go loraDev.SendReceiveLoop()
	}
	// Ref: https://packagewjx.github.io/2018/09/19/cgo-cstring-ram-leak/
	cBufferNeedToFree := C.CBytes(data)
	defer C.free(unsafe.Pointer(cBufferNeedToFree))
	fmt.Printf("Lora.Send: Sending len: %d\n", len(data))
	res := C.send((*C.uchar)(cBufferNeedToFree), C.int(len(data)))
	return int(res)
}

func (loraDev Lora) Receive() []byte {
	fmt.Printf("loraDev: %s, loop started: %v, loraDev: 0x%p\n", loraDev.DeviceName, loraDev.IsHandlingLoopStarted, &loraDev)
	// fmt.Printf("loraDev: %s, %d, %v", loraDev.DeviceName, loraDev.ModuleInst.Freq, loraDev.IsHandlingLoopStarted)
	if C.is_loop_started() == 0 {
		fmt.Printf("Start handling loop\n")
		go loraDev.ReceiveLoop()
	}

	buffer := bytes.NewBuffer(make([]byte, 513))
	log.Println("Got buffer\n")
	// fmt.Println("Got buffer\n")
	// len := 0
	// Call the C function to fill the buffer
	fmt.Printf("receive to: %p\n", &buffer.Bytes()[0])
	len := int(C.receive((*C.uchar)(unsafe.Pointer(&buffer.Bytes()[0])), C.int(buffer.Cap())))
	// len := 0
	log.Println("After calling receive\n")
	// fmt.Println("After calling receive\n")
	if len <= 0 {
		log.Printf("Receive error, len: %d", int(len))
		return make([]byte, 0)
	} else {
		byteSlice := make([]byte, len)
		buffer.Read(byteSlice)
		return byteSlice
	}
}
