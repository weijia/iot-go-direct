package iot_go_test

import (
	"iot_go/pkg/util"
	"testing"
)

func TestFolder(t *testing.T) {
	util.GetAppRoot()
	downloadingFolder := util.GetOrCreateDownloadingFolder()
	util.IotLogInfo(downloadingFolder)
}
