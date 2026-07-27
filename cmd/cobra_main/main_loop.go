package main

import (
	"fmt"
	"os"
	"context"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"iot_go/pkg/util"
	"iot_go/pkg/bsp"
	"iot_go/pkg/main_msg_handler"
	"iot_go/pkg/serial"
)

var loraHost string
var pushHost string
var mock bool
var portName string
var baud int

// versionCmd represents the version command
var mainLoopCmd = &cobra.Command{
	Use:     "main_loop",
	Short:   "Run the main loop",
	Long:    `Run the main loop (serial control board)`,
	Run: func(cmd *cobra.Command, args []string) {
		printVer()
		fmt.Printf("main() pid: %d\n", os.Getpid())

		util.IotLog("lora service host: %s, push host: %s", loraHost, pushHost)
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

		go main_msg_handler.InfiniteAppLoop(ctx, port, mock)

		// 捕捉退出信号，断开连接并退出程序
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()

		bsp.GetBsp().StopAllProcess()

	},
}

func init() {
	pflag.StringVarP(&loraHost, "lora_service_ip",  "l", "127.0.0.1", "Lora service server (legacy)")
	pflag.StringVarP(&pushHost, "push_msg_host", "s", "127.0.0.1", "Lora msg push host server (legacy)")
	pflag.BoolVar(&mock, "mock", false, "use in-memory mock serial transport (no device needed)")
	pflag.StringVar(&portName, "port", "COM1", "serial port, e.g. /dev/ttyS1 or COM1")
	pflag.IntVar(&baud, "baud", 9600, "baud rate")
	// 解析命令行参数
	viper.BindPFlags(pflag.CommandLine)

	rootCmd.AddCommand(mainLoopCmd)
}
