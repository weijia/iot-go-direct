package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"iot_go/pkg/app_manager"
	"iot_go/pkg/msg"
	"iot_go/pkg/serial_gateway"
	"iot_go/pkg/util"
)

// init 复用 cmd/iot/main.go 的"启动更新版本 app"思路：
// 若 apps/ 目录下存在比当前运行的更新版本的 mqtt_serial_gateway 二进制，
// 则 re-exec 拉起该版本（自身热更新），否则继续正常启动。
func init() {
	if latest, err := app_manager.GetLatestApp(msg.FIRMWARE_FOLDER, "mqtt_serial_gateway"); err == nil {
		util.IotLog("Found newer version: %s, re-executing", latest)
		cmd := exec.Command(latest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			os.Exit(0)
		} else {
			util.IotLogErrWithStr("re-exec newer version failed", err)
		}
	}
}

func main() {
	configPath := "mqtt_serial_gateway-prod.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("read config %s: %v", configPath, err)
	}
	var cfg serial_gateway.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse config %s: %v", configPath, err)
	}

	// 合理的默认值（跨平台：Windows 默认 COM3，Linux 默认 /dev/ttyS0）
	if cfg.SerialPort == "" {
		if _, err := os.Stat("/dev"); err == nil {
			cfg.SerialPort = "/dev/ttyS0"
		} else {
			cfg.SerialPort = "COM3"
		}
	}
	if cfg.BaudRate == 0 {
		cfg.BaudRate = 9600
	}
	if cfg.DataBits == 0 {
		cfg.DataBits = 8
	}
	if cfg.StopBits == 0 {
		cfg.StopBits = 1
	}
	if cfg.ClientID == "" {
		cfg.ClientID = fmt.Sprintf("mqtt_serial_gateway-%d", os.Getpid())
	}

	util.ConfigLogFile("mqtt_serial_gateway-log.txt", util.LogConfigParams{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gw, err := serial_gateway.NewGateway(cfg)
	if err != nil {
		log.Fatalf("init gateway: %v", err)
	}
	if err := gw.Run(ctx); err != nil {
		log.Fatalf("run gateway: %v", err)
	}
	util.IotLog("mqtt_serial_gateway started, serial=%s topic(in)=%s topic(out)=%s",
		cfg.SerialPort, cfg.ToSerialTopic, cfg.FromSerialTopic)

	// 等待退出信号，断开连接并退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
	gw.Close()
	util.IotLog("mqtt_serial_gateway stopped")
}
