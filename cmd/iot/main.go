package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
	"iot_go/pkg/main_msg_handler"
	"iot_go/pkg/msg"
	"iot_go/pkg/serial"
	"iot_go/pkg/util"
)

func init() {
	appStorePath := filepath.Join(util.GetAppRoot(), msg.FIRMWARE_FOLDER)
	entries, err := os.ReadDir(appStorePath)
	var latestVersion float64
	var latestApp string
	latestApp = ""

	if err == nil {
		// Iterate through the directory entries and print their names
		for _, entry := range entries {
			filename := entry.Name()
			util.IotLog("Checking app file: %s", filename)
			verStr := strings.Replace(filename, "main-", "", -1)
			verStr = strings.Replace(verStr, ".exe", "", -1)
			version, err := strconv.ParseFloat(verStr, 32)

			if err == nil && version > latestVersion {
				latestVersion = version
				latestApp = filename
			}
		}
		appPath, errGetExecutable := os.Executable()
		curAppName := filepath.Base(appPath)
		latestAppFullPath := filepath.Join(appStorePath, latestApp)
		if errGetExecutable == nil && latestApp != "" && curAppName != latestApp {
			util.IotLog("Starting: %s from curAppName: %s", latestApp, curAppName)
			cmd := exec.Command(latestAppFullPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// Will not return until exit
			err := cmd.Run()
			if err == nil {
				os.Exit(0)
			} else {
				util.IotLogErrWithStr("Error execute app", err)
			}
		}
	} else {
		util.IotLogErrorWithFormatStr("Error reading app store folder: %s with err: %v", appStorePath, err)
	}
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
	var mock bool
	var portName string
	var baud int
	var broker string
	var selftest bool
	flag.BoolVar(&mock, "mock", false, "use in-memory mock serial transport (no device needed)")
	flag.StringVar(&portName, "port", "COM1", "serial port, e.g. /dev/ttyS1 or COM1")
	flag.IntVar(&baud, "baud", 9600, "baud rate")
	flag.StringVar(&broker, "broker", "", "MQTT broker host:port, e.g. broker.emqx.io:1883 (overrides config; empty => bsp default)")
	flag.BoolVar(&selftest, "selftest", false, "after start, internally query the (mock) board and print its reply; does not need MQTT broker")
	flag.Parse()

	util.ConfigLogFile("main-log.txt", bsp.BspConfigInstance.LogConfigParams)
	ctx, cancel := context.WithCancel(context.Background())

	var port *serial.Port
	if mock {
		util.IotLog("Running in MOCK serial mode (no device)")
		port = serial.OpenMock(nil)
	} else {
		p, err := serial.Open(serial.Config{Port: portName, Baud: baud})
		if err != nil {
			util.IotLogErrorWithFormatStr("failed to open serial port: %v", err)
			os.Exit(1)
		}
		port = p
	}

	// 显式指定 -broker 时连接 MQTT(即使 -mock 也连，用于无硬件端到端测试)
	actualSkipMqtt := mock
	if broker != "" {
		bsp.MqttBrokerOverride = broker
		actualSkipMqtt = false
	}
	go main_msg_handler.InfiniteAppLoop(ctx, port, actualSkipMqtt)

	// -selftest: 不依赖 broker，直接在进程内对默认节点发串口查询，
	// 打印 MockBoard(虚拟玻璃)的回包与节点状态回填，用来验证串口回复链路。
	if selftest {
		go func() {
			time.Sleep(2 * time.Second)
			for i := 0; i < 3; i++ {
				frame, ok := controller.Ctrl.QueryStatus()
				if ok {
					zones := serial.ZonesFromPayload(frame.Payload)
					log.Printf("[selftest] MockBoard replied: cmd=%d zones=%v", frame.Cmd, zones)
					controller.Ctrl.UpdateNodeStatesFromSerialReply(frame)
					ns := bsp.GetOrCreateNodeState("node1")
					log.Printf("[selftest] node1 CompletionStatus=%d (0未知/1执行中/2完成)", ns.CompletionStatus)
				} else {
					log.Printf("[selftest] QueryStatus no reply (timeout)")
				}
				time.Sleep(3 * time.Second)
			}
		}()
	}

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()

	bsp.GetBsp().StopAllProcess()
}
