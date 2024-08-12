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
)

var loraHost string
var pushHost string

// versionCmd represents the version command
var mainLoopCmd = &cobra.Command{
	Use:     "main_loop",
	Short:   "Run the main loop without after all lora service started",
	Long:    `Run the main loop without after all lora service started`,
	Run: func(cmd *cobra.Command, args []string) {
		printVer()
		fmt.Printf("main() pid: %d\n", os.Getpid())
		loraServiceIp := loraHost

		util.IotLog("lora service host: %s, push host: %s", loraHost, pushHost)
		// util.ConfigLogFile("main-log.txt", bsp.BspConfigInstance.LogConfigParams)
		ctx, cancel := context.WithCancel(context.Background())
	
		go main_msg_handler.InfiniteAppLoop(ctx, loraServiceIp)
	
		// 捕捉退出信号，断开连接并退出程序
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	
		bsp.GetBsp().StopAllProcess()
		
	},
}

func init() {
	pflag.StringVarP(&loraHost, "lora_service_ip",  "l", "127.0.0.1", "Lora service server")
	pflag.StringVarP(&pushHost, "push_msg_host", "s", "127.0.0.1", "Lora msg push host server")
	// 解析命令行参数
	viper.BindPFlags(pflag.CommandLine)

	rootCmd.AddCommand(mainLoopCmd)
}
