package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/weijia/supervisord/supervisord_main"
	// "iot_go/pkg/bsp"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version number of generated code example",
	Aliases: []string{"start", "status", "stop", "restart", "shutdown", "reload", "signal", "pid", "logtail"},
	Long:    `All software has versions. This is generated code example`,
	Run: func(cmd *cobra.Command, args []string) {

		printVer()
		supervisord_main.RealMain()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVer() {
	fmt.Println("Build Date:", BuildDate)
	fmt.Println("Git Commit:", GitCommit)
	fmt.Println("Version:", Version)
	fmt.Println("Go Version:", GoVersion)
	fmt.Println("OS / Arch:", OsArch)
}
