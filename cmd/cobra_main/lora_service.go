package main

import (
	"github.com/spf13/cobra"
    "github.com/spf13/pflag"
	"iot_go/pkg/lora_rpc"
)

// loraCmd represents the lora service command
var loraCmd = &cobra.Command{
	Use:     "lora",
	Short:   "Start a lora service",
	Aliases: []string{},
	Long:    `All software has versions. This is generated code example`,
	Run: func(cmd *cobra.Command, args []string) {
		printVer()
		// var cmd_name string
		var dev string
		var port int
		var host_ip string

		// pflag.StringVarP(&cmd_name, "cmd", "c","", "The function you want to execute")
		pflag.StringVarP(&dev, "dev", "d", "/dev/spidev1.0", "Path to the device file")
		pflag.IntVarP(&port, "port", "p", 8866, "The port lora service will listen on")
		pflag.StringVarP(&host_ip, "host_ip", "h", "127.0.0.1", "The host ip lora service will send received message to")

		lora_rpc.StartLoraService(dev, port, host_ip)
	},
}

func init() {
	rootCmd.AddCommand(loraCmd)
}
