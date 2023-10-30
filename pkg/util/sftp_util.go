package util

import (
	"fmt"
	"io"
	"iot_go/pkg/shared"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Download(sftpInfo shared.SftpInfo, localPath string, isDownloading bool) error {
	// SSH客户端配置
	config := &ssh.ClientConfig{
		User:            sftpInfo.User,
		Auth:            []ssh.AuthMethod{ssh.Password(sftpInfo.Pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到远程SSH服务器
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sftpInfo.IP, sftpInfo.Port), config)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("无法连接到SSH服务器: %v", err))
		return err
	}
	defer conn.Close()

	// 创建SFTP客户端
	client, err := sftp.NewClient(conn)
	if err != nil {
		IotLogErrorStr(fmt.Sprintf("无法创建SFTP客户端: %v", err))
		return err
	}
	defer client.Close()

	remoteFilePath := sftpInfo.Path
	localFilePath := localPath

	if isDownloading {

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
	} else {
		// 打开本地文件
		fileInfo, err := os.Stat(localFilePath)
		if err != nil {
			IotLogError(err)
			return err
		}

		var srcReader io.Reader

		if fileInfo.IsDir() {
			// Read the contents of the directory
			entries, err := os.ReadDir(localFilePath)
			if err != nil {
				fmt.Println("Error reading the directory:", err)
				return err
			}
			var filenameList string

			// Iterate through the directory entries and print their names
			for _, entry := range entries {
				filenameList += entry.Name() + "\n"
			}
			srcReader = strings.NewReader(filenameList)
		} else {
			localFile, err := os.Open(localFilePath)
			if err != nil {
				IotLogError(err)
				return err
			}
			defer localFile.Close()
			srcReader = localFile
		}

		// 创建远程文件
		remoteFile, err := client.Create("/home/ota/upgrade/upload")
		if err != nil {
			IotLogError(err)
			return err
		}
		defer remoteFile.Close()

		// 将本地文件内容复制到远程文件
		bytes, err := io.Copy(remoteFile, srcReader)
		if err != nil {
			IotLogError(err)
			return err
		}

		IotLog("Successfully uploaded %d bytes.\n", bytes)
	}

	return nil
}
