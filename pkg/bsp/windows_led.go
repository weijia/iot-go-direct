//go:build windows

package bsp

import "iot_go/pkg/util"



func WriteLedFile(devName string, value int) {

	util.IotLog("Write to LED device: %s, value: %d", devName, value)
}