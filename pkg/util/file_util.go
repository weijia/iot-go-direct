package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func PruneFolderSoRemainingFileNumberEqualTo(folderPath string, remainingFileNum int) {
	// List files in the folder
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("Error reading folder: %v", err))
		return
	}

	// Create a slice to store file info
	var fileInfoSlice []os.FileInfo

	fileInfoSlice = append(fileInfoSlice, files...)

	// Sort files by modification time in ascending order
	sort.Slice(fileInfoSlice, func(i, j int) bool {
		return fileInfoSlice[i].ModTime().Before(fileInfoSlice[j].ModTime())
	})

	// Delete the oldest files
	deleteFileNum := len(fileInfoSlice) - remainingFileNum
	deleteCount := 0
	for _, fileInfo := range fileInfoSlice {
		filePath := folderPath + "/" + fileInfo.Name()
		err := os.Remove(filePath)
		if err != nil {
			IotLogErrorStr(fmt.Sprintf("Error deleting file %s: %v", fileInfo.Name(), err))
		} else {
			IotLogErrorStr(fmt.Sprintf("Deleted file: %s\n", fileInfo.Name()))
			deleteCount++
		}

		// Stop after deleting 5 files
		if deleteCount >= deleteFileNum {
			break
		}
	}
}