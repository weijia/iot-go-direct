package app_manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var readDirFunc = os.ReadDir

func GetLatestApp(appFolder string, executableName string) (string, error) {

	appStorePath := filepath.Join(GetAppRoot(), appFolder)
	entries, err := readDirFunc(appStorePath)
	var latestVersion float64
	var latestApp string
	latestApp = ""

	if err == nil {
		// Iterate through the directory entries and print their names
		for _, entry := range entries {
			filename := entry.Name()
			fmt.Printf("Checking app file: %s\n", filename)
			appBaseNamePrefix := fmt.Sprintf("%s-", executableName)
			verStr := strings.Replace(filename, appBaseNamePrefix, "", -1)
			verStr = strings.Replace(verStr, ".exe", "", -1)
			version, err := strconv.ParseFloat(verStr, 32)

			if err == nil && version > latestVersion {
				latestVersion = version
				latestApp = filename
			}
		}
		appPath, errGetExecutable := GetExecutable()
		
		if errGetExecutable == nil && latestApp != "" {
			curAppName := filepath.Base(appPath) // filename of the app
			if curAppName != latestApp {
				latestAppFullPath := filepath.Join(appStorePath, latestApp)
				fmt.Printf("Latest app: %s from curAppName: %s\n", latestApp, curAppName)
				return latestAppFullPath, nil
			} else {
				return "", os.ErrNotExist
			}

		}
	} 
	fmt.Printf("Error reading app store folder: %s with err: %v\n", appStorePath, os.ErrNotExist)
	return "", os.ErrNotExist
}
