package util

import (
	"os"
	"path/filepath"
	"syscall"
)

// Ref: https://blog.csdn.net/qq_33747476/article/details/126056458
func fork1() (*os.Process, error) {

	// 获取当前进程名和工作目录
	execName, err := os.Executable() // 输出与 os.Args[0] 类似
	if err != nil {
		return nil, err
	}
	// 生成子进程
	p, err := os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir: filepath.Dir(execName), //工作路径 不包含文件名
		Env: os.Environ(),           // 获取当前环境变量 并且传入子进程
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: &syscall.SysProcAttr{},
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}
