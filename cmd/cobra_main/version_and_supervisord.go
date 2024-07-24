package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"iot_go/version"
	"github.com/weijia/supervisord/supervisord_main"
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
	fmt.Println("Build Date:", version.BuildDate)
	fmt.Println("Git Commit:", version.GitCommit)
	fmt.Println("Version:", version.Version)
	fmt.Println("Go Version:", version.GoVersion)
	fmt.Println("OS / Arch:", version.OsArch)
}
