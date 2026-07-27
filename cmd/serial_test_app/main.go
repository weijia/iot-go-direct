package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
	"iot_go/pkg/serial"
)

// 测试 app：完全不接串口设备，用 serial.OpenMock 的内存管道 + 模拟板
// 验证帧编解码、node↔zone 映射、MsgLoop 半双工、上电中超时等逻辑。

func check(cond bool, name string) {
	if cond {
		fmt.Printf("  [PASS] %s\n", name)
	} else {
		fmt.Printf("  [FAIL] %s\n", name)
	}
}

func zoneStatusText(s byte) string {
	switch s {
	case serial.ZoneCompleted:
		return "A完成"
	case serial.ZoneExecuting:
		return "B执行中"
	case serial.ZonePowering:
		return "C上电中"
	default:
		return "F保留"
	}
}

func printZones(zones [16]byte) {
	for i := 0; i < 16; i++ {
		fmt.Printf("  区%2d: %s\n", i+1, zoneStatusText(zones[i]&0xF))
	}
}

func runSelfTest(c *controller.Controller, boardPowering *bool) {
	fmt.Println("==== 自动化自测开始 ====")

	// 1) node1 改色：映射验证
	node1Zones, ok := controller.NodeColorToZones("node1", "1022344052647082")
	check(ok, "NodeColorToZones(node1) 成功")
	req, ok := func() (serial.Frame, bool) {
		frame := serial.BuildFrame(serial.CmdChangeColor, serial.PayloadFromZones(node1Zones))
		return serial.ParseFrame(frame[:])
	}()
	check(ok, "node1 改色请求帧可解析")
	if ok {
		rz := serial.ZonesFromPayload(req.Payload)
		expect := [8]byte{0, 2, 4, 0, 2, 4, 0, 2}
		mappedOK := true
		for i := 0; i < 8; i++ {
			if rz[i] != expect[i] {
				mappedOK = false
			}
		}
		for i := 8; i < 16; i++ {
			if rz[i] != serial.ColorKeep {
				mappedOK = false
			}
		}
		check(mappedOK, "node1 映射: zone1~8 着色 / zone9~16 保留")
	}

	// 2) node2 改色：映射验证
	node2Zones, ok := controller.NodeColorToZones("node2", "1022344052647082")
	check(ok, "NodeColorToZones(node2) 成功")
	req2, ok := func() (serial.Frame, bool) {
		frame := serial.BuildFrame(serial.CmdChangeColor, serial.PayloadFromZones(node2Zones))
		return serial.ParseFrame(frame[:])
	}()
	check(ok, "node2 改色请求帧可解析")
	if ok {
		rz := serial.ZonesFromPayload(req2.Payload)
		mappedOK := true
		for i := 0; i < 8; i++ {
			if rz[i] != serial.ColorKeep {
				mappedOK = false
			}
		}
		expect := [8]byte{0, 2, 4, 0, 2, 4, 0, 2}
		for i := 8; i < 16; i++ {
			if rz[i] != expect[i-8] {
				mappedOK = false
			}
		}
		check(mappedOK, "node2 映射: zone1~8 保留 / zone9~16 着色")
	}

	// 3) 正常改色回包 + 状态回填（handler 内会调用，这里手动调以验证映射）
	frame3, ok := c.ChangeColorForZones(node1Zones)
	check(ok, "node1 改色收到板子应答")
	if ok {
		c.UpdateNodeStatesFromSerialReply(frame3)
	}
	check(bsp.GetOrCreateNodeState("node1").CompletionStatus == 1, "node1 改色后 CompletionStatus=执行中")

	// 4) 查询状态回包
	f, ok := c.QueryStatus()
	check(ok && f.Cmd == serial.StatusQueryStatus, "查询状态收到板子应答")
	if ok {
		printZones(serial.ZonesFromPayload(f.Payload))
	}

	// 5) 上电中阻塞 -> 超时
	*boardPowering = true
	start := time.Now()
	_, ok = c.ChangeColorForZones(node1Zones)
	elapsed := time.Since(start)
	check(!ok && elapsed >= controller.ReplyTimeout-200*time.Millisecond,
		fmt.Sprintf("上电中改色超时无应答 (耗时 %v)", elapsed.Round(time.Millisecond)))
	*boardPowering = false

	// 6) 退出上电中后查询正常
	_, ok = c.QueryStatus()
	check(ok, "退出上电中后查询状态恢复正常")

	fmt.Println("==== 自动化自测结束 ====")
}

func interactive(c *controller.Controller, boardPowering *bool) {
	fmt.Println("\n==== 进入交互模式 (命令: color <nodeId> <hex> | status | power on|off | quit) ====")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "color":
			if len(parts) < 3 {
				fmt.Println("用法: color <nodeId> <hex如 1022344052647082>")
				continue
			}
			zones, ok := controller.NodeColorToZones(parts[1], parts[2])
			if !ok {
				fmt.Println("nodeId 非法或颜色串错误")
				continue
			}
			f, ok := c.ChangeColorForZones(zones)
			if !ok {
				fmt.Println("改色超时无应答(可能处于上电中)")
				continue
			}
			fmt.Printf("改色应答 状态码=%d\n", f.Cmd)
		case "status":
			f, ok := c.QueryStatus()
			if !ok {
				fmt.Println("查询超时无应答")
				continue
			}
			fmt.Printf("状态应答 状态码=%d\n", f.Cmd)
			printZones(serial.ZonesFromPayload(f.Payload))
		case "power":
			if len(parts) >= 2 && parts[1] == "on" {
				*boardPowering = true
				fmt.Println("板子进入上电中(丢弃所有帧)")
			} else {
				*boardPowering = false
				fmt.Println("板子退出上电中")
			}
		case "quit", "exit":
			return
		default:
			fmt.Println("未知命令")
		}
	}
}

func main() {
	// 为测试设置节点列表，使 node↔zone 映射可用（无需真实配置文件）。
	bsp.BspConfigInstance.BaseConfigParams.NodeList1 = []string{"node1"}
	bsp.BspConfigInstance.BaseConfigParams.NodeList2 = []string{"node2"}

	// 用内存 mock 控制板替代真实串口。
	boardPowering := false
	port := serial.OpenMock(&boardPowering)
	c := controller.Init(port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.MsgLoop(ctx)

	// 模拟主循环：读协程的帧 -> 转发给 controller 等待者（沿用真实架构约定）。
	go func() {
		for f := range port.RecvCh() {
			c.RecvCh() <- f
		}
	}()

	runSelfTest(c, &boardPowering)
	interactive(c, &boardPowering)
}
