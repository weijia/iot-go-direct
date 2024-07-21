package main

import (
	"iot_go/pkg/lora_rpc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// var cmd_name string
var dev string
var port int
var host_ip string

// loraCmd represents the lora service command
var loraCmd = &cobra.Command{
	Use:     "lora",
	Short:   "Start a lora service",
	Aliases: []string{},
	Long:    `All software has versions. This is generated code example`,
	Run: func(cmd *cobra.Command, args []string) {
		printVer()

		lora_rpc.StartLoraService(dev, port, host_ip)
	},
}

func init() {
	// pflag.StringVarP(&cmd_name, "cmd", "c","", "The function you want to execute")
	pflag.StringVarP(&dev, "dev", "d", "/dev/spidev1.0", "Path to the device file")
	pflag.IntVarP(&port, "port", "p", 8866, "The port lora service will listen on")
	pflag.StringVarP(&host_ip, "host_ip", "i", "127.0.0.1", "The host ip lora service will send received message to")

	// 解析命令行参数
	viper.BindPFlags(pflag.CommandLine)
	rootCmd.AddCommand(loraCmd)
}
