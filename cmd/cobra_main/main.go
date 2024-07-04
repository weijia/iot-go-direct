package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "MyApp is a sample CLI application",
		Long:  `MyApp is a sample command line interface (CLI) application using Cobra.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Main command logic here
			fmt.Println("Hello from MyApp!")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
