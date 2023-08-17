package iot_go_test



import (
	"encoding/json"
	"fmt"
	"testing"

	"iot_go/pkg/util"
    "iot_go/pkg/shared"
)

func TestNodeInfoMsg(t *testing.T) {
    var sftpInfo shared.SftpInfo
    sftpInfo.User = "test"
    sftpInfo.Pwd = "pwd"
    sftpInfo.IP = "ipaddress"
    sftpInfo.Port = 22
    sftpInfo.Path = "/path/to/file"
	util.Download(sftpInfo)
}
