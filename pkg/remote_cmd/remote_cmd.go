package remote_cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type RemoteCmdResult struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}

func ExecuteCmd(cmdStr string) RemoteCmdResult {
	var res RemoteCmdResult
	// 定义要执行的命令
	// 注意：这里使用了一个可能会产生错误的命令作为示例
	// 假设我们有一个包含命令和参数的字符串

	// 使用strings.Fields来分割字符串为命令和参数列表
	// 注意：这个方法会按照空格分割字符串，所以如果你的参数中包含空格，你需要使用引号或其他方法来避免被错误分割
	parts := strings.Fields(cmdStr)

	// 第一个元素是命令名，其余的是参数
	cmdName := parts[0]
	cmdArgs := parts[1:]

	// 使用exec.Command执行命令
	cmd := exec.Command(cmdName, cmdArgs...)

	// 创建两个buffer分别用于存储标准输出和标准错误输出
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// 将命令的标准输出和标准错误输出分别设置到两个buffer中
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		// 如果命令执行出错，这里会捕获到错误，但标准错误输出已经存储在stderr buffer中
		fmt.Println("执行命令出错:", err)
	}

	// 打印命令的标准输出
	fmt.Println("标准输出:", stdout.String())

	// 打印命令的标准错误输出
	// 注意：即使命令执行出错，这里的打印仍然会执行
	fmt.Println("标准错误输出:", stderr.String())

	res.StdErr = stderr.String()
	res.StdOut = stdout.String()
	return res

}
