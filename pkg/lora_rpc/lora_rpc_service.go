package lora_rpc

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/lora"
	"iot_go/pkg/lora_shared"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

type LoraRpc struct {
	Port    int
	LoraDev *lora.Lora
}

func (loraRpc LoraRpc) InitLora(argType shared.Module, reply *lora_shared.ReplyResult) error {
	reply.Result = loraRpc.LoraDev.InitLora(argType)
	return nil
}

func (loraRpc LoraRpc) Exit(argType lora_shared.EmptyArg, reply *lora_shared.ReplyResult) error {
	reply.Result = loraRpc.LoraDev.Exit()
	return nil
}

func (loraRpc LoraRpc) Send(argType lora_shared.LoraData, reply *lora_shared.ReplyResult) error {
	// fmt.Println("RPC: Send called")
	// log.Println("RPC: Send called")
	fmt.Printf("Sending len: %d\n", len(argType.Data))
	reply.Result = loraRpc.LoraDev.Send(argType.Data)
	return nil
}

// func (loraRpc LoraRpc) Receive(argType lora_shared.EmptyArg, reply *lora_shared.LoraData) error {
// 	fmt.Printf("RPC: Receive called, Rpc: %p, dev: %p\n", &loraRpc, &loraRpc.LoraDev)
// 	log.Println("RPC: Receive called")
// 	reply.Data = loraRpc.LoraDev.Receive()
// 	return nil
// }

// func (loraRpc LoraRpc) ToggleDebug(argType lora_shared.EmptyArg, reply *lora_shared.EmptyArg) error {
// 	loraRpc.LoraDev.ToggleDebug()
// 	return nil
// }

func NewLoraRpc(devName string, port int) *LoraRpc {
	loraDev := lora.NewLora(devName)
	return &LoraRpc{
		Port:    port,
		LoraDev: loraDev,
	}
}

var ModuleIndex int

func GetModuleIndexFromDevName(devName string) int {
	moduleIndex := 0
	dev := strings.Split(devName, "/")
	res := strings.ReplaceAll(dev[2], "spidev", "")
	devNumFloat , err := strconv.ParseFloat(res, 32)
	if err != nil {
		util.IotLogErrWithStr("Parse spi dev num error", err)
	} else {
		moduleIndex = int(devNumFloat)
	}
	return moduleIndex
}

func StartLoraServiceInBackground(devName string, port int, pushHost string) {
	moduleIndex := GetModuleIndexFromDevName(devName)

	util.ConfigLogFile(fmt.Sprintf("%d-log.txt", moduleIndex), 
			bsp.BspConfigInstance.LogConfigParams)

	// file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	util.IotLogError(err)
	// }
	// defer file.Close()

	// 设置日志输出到文件
	// log.SetOutput(file)
	StartLoraService(devName, port, pushHost)
}

func StartLoraService(devName string, port int, pushHost string) {
	bsp.InitConfig()

	util.IotLogInfo(fmt.Sprintf("Starting lora service on dev: %s, port: %d\n", devName, port))
	go lora.PushLoraMsgToRpc(8869, GetModuleIndexFromDevName(devName), pushHost)

	loraRpc := NewLoraRpc(devName, port)
	// The loraRpc object will be copied to RPC procedure instead of sending the original object
	// So the data changed after Register may be discarded
	rpc.Register(loraRpc)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("serve error: %v", err))
	}
}
