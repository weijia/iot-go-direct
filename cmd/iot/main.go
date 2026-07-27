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

			// 验证按协议格式的改色请求(node_id=12345678, 8-hex)能映射到区并触发虚拟玻璃回包。
			// 这正是之前被默认节点列表(node1/node2 占位符)挡掉、回 invalid 色的场景。
			color := "1222324252627282"
			zones, ok := controller.NodeColorToZones("12345678", color)
			if !ok {
				log.Printf("[selftest] NodeColorToZones(12345678) FAILED: node not in NodeList1")
				return
			}
			frame, sent := controller.Ctrl.ChangeColorForZones(zones)
			if !sent || frame.Cmd != serial.StatusChangeColor {
				log.Printf("[selftest] update_glass_color 12345678: board no reply (timeout)")
				return
			}
			controller.Ctrl.UpdateNodeStatesFromSerialReply(frame)
			ns := bsp.GetOrCreateNodeState("12345678")
			log.Printf("[selftest] update_glass_color 12345678 OK: reply cmd=%d, node CompletionStatus=%d, requested color=%s",
				frame.Cmd, ns.CompletionStatus, color)

			// 等待控制板完成变色(mock 约 800ms)，并让主循环自动轮询(每 10s)刷新状态，
			// 验证不再需要手动 get_glass_status，心跳即可反映“已完成(completion_status=2)”。
			time.Sleep(10 * time.Second)
			ns2 := bsp.GetOrCreateNodeState("12345678")
			log.Printf("[selftest] after auto-poll: 12345678 CompletionStatus=%d (expect 2=已完成)", ns2.CompletionStatus)
		}()
	}

	// 捕捉退出信号，断开连接并退出程序
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()

	bsp.GetBsp().StopAllProcess()
}
