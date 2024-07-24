//go:build linux
package bsp


import (
	"path/filepath"
	"os"
	"fmt"
)

func WriteLedFile(devName string, value int) {
	rootPath := "/sys/devices/platform/leds/leds"
	devPath := filepath.Join(filepath.Join(rootPath, devName), "brightness")
	// 打开设备文件以进行写入操作
	file, err := os.OpenFile(devPath, os.O_WRONLY, 0)
	if err != nil {
		// 处理错误
		panic(err)
	}
	defer file.Close()

	// 要写入的数据
	data := []byte(fmt.Sprintf("%d", value))

	// 将数据写入设备文件
	_, err = file.Write(data)
	if err != nil {
		// 处理写入错误
		panic(err)
	}
	// util.IotLog("Write to LED device: %s, value: %d", devPath, value)
}