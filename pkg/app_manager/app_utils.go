package app_manager

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetExecutable() (string, error) {
	// Get the current executable full path
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return "", err
	}

	executablePath, err := filepath.EvalSymlinks(executable)
	if err != nil {
		fmt.Println("Error resolving symlinks:", err)
		return "", err
	}

	fmt.Println("Current executable path:", executablePath)
	return executablePath, nil
}
