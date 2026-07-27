package controller

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"iot_go/pkg/bsp"
	"iot_go/pkg/serial"
)

// 16 区与现有 LoRa 节点的对应：
//   node1 (NodeList1) -> 1~8 区  (zone 索引 0..7)
//   node2 (NodeList2) -> 9~16 区 (zone 索引 8..15)
func NodeToZoneRange(nodeId string) (start, end int, ok bool) {
	if bsp.IsInNodeList1(nodeId) {
		return 0, 8, true
	}
	if bsp.IsInNodeList2(nodeId) {
		return 8, 16, true
	}
	return 0, 0, false
}

// NodeColorToZones 把某节点 8 区颜色串转成完整 16 区数组。
// 属于另一节点的区填 ColorKeep(0xF)，这样一条改色帧不会误改另一节点。
//
// colorStr 格式：十六进制对 "区(1-8)颜色"，如 "1f2f3f4f5f6f7f8f"。
func NodeColorToZones(nodeId string, colorStr string) ([16]byte, bool) {
	start, end, ok := NodeToZoneRange(nodeId)
	if !ok {
		return [16]byte{}, false
	}
	var zones [16]byte
	for i := range zones {
		zones[i] = serial.ColorKeep
	}
	data, err := hex.DecodeString(colorStr)
	if err != nil {
		return [16]byte{}, false
	}
	for _, b := range data {
		area := int((b >> 4) & 0xF) // 1..8
		color := b & 0xF
		if area < 1 || area > 8 {
			continue
		}
		idx := start + (area - 1)
		if idx < end {
			zones[idx] = color
		}
	}
	return zones, true
}

// ZonesToNodeColor 是 NodeColorToZones 的逆操作：把某节点的 8 区抽回
// "区 颜色" 十六进制串格式。
func ZonesToNodeColor(nodeId string, zones [16]byte) string {
	start, end, ok := NodeToZoneRange(nodeId)
	if !ok {
		return ""
	}
	res := ""
	for i := start; i < end; i++ {
		area := i - start + 1
		res += fmt.Sprintf("%x%x", area, zones[i]&0xF)
	}
	return res
}

// ZoneStatusToCompletion 把分区状态 nibble 映射到 NodeState 里的
// CompletionStatus 语义（0=错误/进行中, 2=已完成）。
func ZoneStatusToCompletion(status byte) int {
	switch status {
	case serial.ZoneCompleted:
		return 2
	case serial.ZoneExecuting, serial.ZonePowering:
		return 1
	default: // F 保留 / 未知
		return 0
	}
}

// AllZonesCompleted 报告 [start,end) 内所有区是否均为 A(已完成)。
func AllZonesCompleted(zones [16]byte, start, end int) bool {
	for i := start; i < end; i++ {
		if zones[i] != serial.ZoneCompleted {
			return false
		}
	}
	return true
}

// ParseColorNibble 从颜色串里取第一个十六进制位作为区颜色 nibble。
// 串口协议里广播/群组用单一颜色控制所有区；具体格式待定(设置最后再定)。
func ParseColorNibble(s string) byte {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return serial.ColorKeep
	}
	if v, err := strconv.ParseUint(s[:1], 16, 8); err == nil {
		return byte(v) & 0xF
	}
	return serial.ColorKeep
}
