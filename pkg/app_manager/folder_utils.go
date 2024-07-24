package app_manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DOWNLOADING_FOLDER = "downloading"
	FIRMWARE_FOLDER    = "apps"
)

func GetAppRoot() string {
	appPath, err := os.Executable()
	fmt.Printf("appPath: %s\n", appPath)
	if err != nil {
		fmt.Println("Getting executable with err:", err)
		return "./"
	}
	executableDir := filepath.Dir(appPath)
	fmt.Printf("Root path: %s\n", executableDir)
	appRoot := executableDir
	if strings.Contains(executableDir, FIRMWARE_FOLDER) {
		appRoot = filepath.Dir(executableDir)
	}
	return appRoot
}

func GetOrCreateDownloadingFolder() string {
	appRoot := GetAppRoot()
	downloadingFolder := filepath.Join(appRoot, DOWNLOADING_FOLDER)
	os.MkdirAll(downloadingFolder, 0771)
	return downloadingFolder
}

func GetOrCreateAppFolder() string {
	appRoot := GetAppRoot()
	appFolder := filepath.Join(appRoot, FIRMWARE_FOLDER)
	os.MkdirAll(appFolder, 0771)
	return appFolder
}
