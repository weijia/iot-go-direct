package util

import (
	"os"
	"strings"
	"path/filepath"
)

const (
	DOWNLOADING_FOLDER = "downloading"
	FIRMWARE_FOLDER    = "apps"
	CONFIG_FILE_NAME = "iot_go.json"
)

func GetAppRoot() string {
	appPath, err := os.Executable()
	IotLog("appPath: %s", appPath)
	if err != nil {
		IotLogErrWithStr("Getting executable with err:", err)
		return "./"
	}
	executableDir := filepath.Dir(appPath)
	IotLog("%s", executableDir)
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
