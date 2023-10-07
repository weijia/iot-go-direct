package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"

	"iot_go/pkg/bsp"
	"iot_go/pkg/lora_rpc"
	"iot_go/pkg/main_msg_handler"
	"iot_go/pkg/msg"
	"iot_go/pkg/util"

	"github.com/kraken-hpc/go-fork"
)

func init() {
	entries, err := os.ReadDir(msg.FIRMWARE_FOLDER)
	var latestVersion float64
	var latestApp string
	latestApp = ""

	if err == nil {
		// Iterate through the directory entries and print their names
		for _, entry := range entries {
			filename := entry.Name()
			fmt.Println(filename)
			verStr := strings.Replace(filename, "main-", "", -1)
			verStr = strings.Replace(verStr, ".exe", "", -1)
			version, err := strconv.ParseFloat(verStr, 32)

			if err == nil && version > latestVersion {
				latestVersion = version
				latestApp = filename
			}
		}
		if latestApp != "" {
			cmd := exec.Command(filepath.Join(msg.FIRMWARE_FOLDER, latestApp))
			// Will not return until exit
			err := cmd.Run()
			if err == nil {
				os.Exit(0)
			} else {
				util.IotLogErrWithStr("Error execute app", err)
			}
		}
	} else {
		util.IotLogErrWithStr("Error reading the directory", err)
	}

	fork.RegisterFunc("StartLoraService", lora_rpc.StartLoraServiceInBackground)
	fork.RegisterFunc("StartLoraService1", lora_rpc.StartLoraServiceInBackground)
	fork.RegisterFunc("StartLoraService2", lora_rpc.StartLoraServiceInBackground)
	fork.Init()
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
	var loraHost string
	var pushHost string
	flag.StringVar(&loraHost, "s", "127.0.0.1", "Lora service server")
	flag.StringVar(&pushHost, "p", "127.0.0.1", "Lora msg push host server")
	loraServiceIp := loraHost
	flag.Parse()
	// loraServiceIp := "192.168.1.20"
	gatewayPushMsgIp := pushHost
	// gatewayPushMsgIp := "192.168.1.18"
	util.IotLog("lora service host: %s, push host: %s", loraHost, pushHost)

	if err := fork.Fork("StartLoraService", "/dev/spidev1.0", 8866, gatewayPushMsgIp); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("failed to fork: %v", err))
	}
	if err := fork.Fork("StartLoraService1", "/dev/spidev2.0", 8867, gatewayPushMsgIp); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("failed to fork: %v", err))
	}
	if err := fork.Fork("StartLoraService2", "/dev/spidev3.0", 8868, gatewayPushMsgIp); err != nil {
		util.IotLogErrorStr(fmt.Sprintf("failed to fork: %v", err))
	}

	util.ConfigLogFile("main-log.txt", bsp.BspConfigInstance.LogConfigParams)
	ctx, cancel := context.WithCancel(context.Background())

	go main_msg_handler.InfiniteAppLoop(ctx, loraServiceIp)

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()

	bsp.GetBsp().StopAllProcess()
}
