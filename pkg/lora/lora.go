package lora

import (
	"fmt"
	"iot_go/pkg/bsp"
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

func StartLoraService(devName string, port int) {
	bsp.InitConfig()
	// dev := strings.Split(devName, "/")

	// file, err := os.OpenFile("e:\\wwj\\codes\\go\\iot_go\\"+dev[2]+"log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// // 设置日志输出到文件
	// log.SetOutput(file)

	// 以下是一些日志示例
	// log.Println("这是一条普通日志")
	// log.Printf("这是一条带参数的日志：%s", "参数值")
	// log.Fatalf("发生了严重错误：%s", "错误信息")

	fmt.Printf("Starting lora service on dev: %s, port: %d\n", devName, port)
	rolaDev := NewLora(devName)
	rpc.Register(rolaDev)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("serve error:", err)
	}
}
