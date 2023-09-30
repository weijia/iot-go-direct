package util

import (
	"fmt"
	"iot_go/pkg/shared"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Download(sftpInfo shared.SftpInfo, targetPath string) error {
	// SSH客户端配置
	config := &ssh.ClientConfig{
		User:            sftpInfo.User,
		Auth:            []ssh.AuthMethod{ssh.Password(sftpInfo.Pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到远程SSH服务器
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sftpInfo.IP, sftpInfo.Port), config)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("无法连接到SSH服务器：%v", err))
		return err
	}
	defer conn.Close()

	// 创建SFTP客户端
	client, err := sftp.NewClient(conn)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("无法创建SFTP客户端：%v", err))
		return err
	}
	defer client.Close()

	remoteFilePath := sftpInfo.Path
	localFilePath := targetPath

	// 打开远程文件
	remoteFile, err := client.Open(remoteFilePath)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("无法打开远程文件：%v", err))
		return err
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(localFilePath)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("无法创建本地文件：%v", err))
		return err
	}
	defer localFile.Close()

	// 将远程文件内容拷贝到本地文件
	bytes, err := remoteFile.WriteTo(localFile)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("文件拷贝失败：%v", err))
		return err
	}

	fmt.Printf("文件下载成功，共写入%d字节\n", bytes)
	return nil
}
