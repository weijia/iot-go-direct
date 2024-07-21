package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/weijia/supervisord/supervisord_main"

	"iot_go/pkg/app_manager"
	"iot_go/pkg/msg"
	"iot_go/pkg/supervisord_manager"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "MyApp is a sample CLI application",
	Long:  `MyApp is a sample command line interface (CLI) application using Cobra.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main command logic here
		// Create supervisord_config file

		// 创建一个空的字符串到字符串的映射  
		stringMap := make(map[string]string)

		appFullPath, err := app_manager.GetLatestApp(msg.FIRMWARE_FOLDER, "main")

		if err != nil {
			// No updates, so do not need to run main
			curRunningExecutable, getExecutableErr := app_manager.GetExecutable()
			if getExecutableErr != nil {
				log.Fatal("Can not get executable")
			}
			appFullPath = curRunningExecutable
		}
		stringMap["main"] = fmt.Sprintf("%s main_loop", strings.Replace(appFullPath, "\\", "\\\\", -1))


		stringMap["lora1"] = fmt.Sprintf("%s lora --dev /dev/spidev1.0 --port 8866 --host_ip 127.0.0.1", strings.Replace(appFullPath, "\\", "\\\\", -1))
		stringMap["lora2"] = fmt.Sprintf("%s lora --dev /dev/spidev2.0 --port 8867 --host_ip 127.0.0.1", strings.Replace(appFullPath, "\\", "\\\\", -1))
		stringMap["lora3"] = fmt.Sprintf("%s lora --dev /dev/spidev3.0 --port 8868 --host_ip 127.0.0.1", strings.Replace(appFullPath, "\\", "\\\\", -1))
		
		supervisord_manager.GenerateConfig(stringMap)

		supervisord_main.RealMain()
		
	},
}

func main() {
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
