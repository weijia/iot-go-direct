package supervisord_manager

import (
	"fmt"
	"os"
)

func GenerateConfig(appNameToCmdMapping map[string]string) {
	// Create supervisord_config file
	// Define the filename and the multi-line string content
	content := `[inet_http_server]
port=0.0.0.0:9001

`
	outputToStdOut := ",/dev/stdout"
	for name, cmd := range appNameToCmdMapping {
		content = content + fmt.Sprintf(`[program:%s]
command = %s
stdout_logfile = %s.log%s
stdout_logfile_maxbytes = 655350
stdout_logfile_backups = 2
stderr_logfile = %s-err.log
stderr_logfile_maxbytes = 65535
stderr_logfile_backups = 2
autorestart = true
`, name, cmd, name, outputToStdOut, name)
		outputToStdOut = ""
	}

	// Create the file
	filePath := "supervisord.ini"

	// 尝试以只写模式打开文件，如果文件不存在则创建它
	// 使用 os.O_CREATE | os.O_WRONLY | os.O_EXCL 以确保文件不存在时创建
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		// 处理其他类型的错误
		fmt.Println("创建文件时发生错误:", err)
		return
	}
	defer file.Close()

	// Write the multi-line string content to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File created and content written successfully.")
}
