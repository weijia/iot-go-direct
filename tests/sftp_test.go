package iot_go_test

import (
	"testing"

	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

func TestDownload(t *testing.T) {
	var sftpInfo shared.SftpInfo
	sftpInfo.User = "ubuntu"
	sftpInfo.Pwd = "123456Rich"
	sftpInfo.IP = "115.159.53.168"
	sftpInfo.Port = 22
	sftpInfo.Path = "/path/to/file"
	util.Download(sftpInfo)
}
