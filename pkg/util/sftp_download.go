package util

import (
	"fmt"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Download(sftpInfo shared.SftpInfo) {
	// SSH客户端配置
	config := &ssh.ClientConfig{
		User:            sftpInfo.User,
		Auth:            []ssh.AuthMethod{ssh.Password(sftpInfo.Pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到远程SSH服务器
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sftpInfo.IP, sftpInfo.Port), config)
	if err != nil {
		util.IotLogErrorStr(fmt.Sprintf("无法连接到SSH服务器：%v", err))
	}
	defer conn.Close()

	// 创建SFTP客户端
	client, err := sftp.NewClient(conn)
	if err != nil {
		util.IotLogErrorStr(fmt.Sprintf("无法创建SFTP客户端：%v", err))
	}
	defer client.Close()

	remoteFilePath := sftpInfo.Path
	localFilePath := "firmware/" + filepath.Base(sftpInfo.Path)

	// 打开远程文件
	remoteFile, err := client.Open(remoteFilePath)
	if err != nil {
		util.IotLogErrorStr(fmt.Sprintf("无法打开远程文件：%v", err))
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(localFilePath)
	if err != nil {
		util.IotLogErrorStr(fmt.Sprintf("无法创建本地文件：%v", err))
	}
	defer localFile.Close()

	// 将远程文件内容拷贝到本地文件
	bytes, err := remoteFile.WriteTo(localFile)
	if err != nil {
		util.IotLogErrorStr(fmt.Sprintf("文件拷贝失败：%v", err))
	}

	fmt.Printf("文件下载成功，共写入%d字节\n", bytes)
}
