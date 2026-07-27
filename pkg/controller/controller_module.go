package controller

import (
	"context"
	"time"

	"iot_go/pkg/bsp"
	"iot_go/pkg/serial"
	"iot_go/pkg/util"
)

const defaultReplyTimeout = 3 * time.Second

// ReplyTimeout 是等待板子应答的超时，测试 app 与业务层可引用。
var ReplyTimeout = defaultReplyTimeout

// Controller 是单串口控制板的抽象，替代原 lora_module(3 模块) + lora_rpc。
type Controller struct {
	port    *serial.Port
	timeout time.Duration
	// SendingCh 承载业务层发来的发送请求，由 MsgLoop 串行消费。
	SendingCh chan sendReq
	// recvCh 是主循环在更新节点状态后，把板子应答转发过来的通道，
	// MsgLoop 里的等待者(sendAndWait)从这里取回复。
	// 特意用无缓冲：若应答在请求已超时后才到，转发会被丢弃，
	// 不会被误认成下一次请求的回复。
	recvCh chan serial.Frame
}

type sendReq struct {
	frame   [serial.FrameLen]byte
	replyCh chan serial.Frame
}

// Ctrl 是进程级单例控制板实例。
var Ctrl *Controller

func Init(port *serial.Port) *Controller {
	c := &Controller{
		port:      port,
		timeout:   defaultReplyTimeout,
		SendingCh: make(chan sendReq, 16),
		recvCh:    make(chan serial.Frame),
	}
	Ctrl = c
	return c
}

// MsgLoop 串行化所有板子通信：一次只发一帧，然后等主循环转发的回复。
// 这样既保证半双工单在途，又复用现有"应用数据只在主循环改"的约定。
func (c *Controller) MsgLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			util.IotLogInfo("Controller.MsgLoop: context cancelled, quit")
			return
		case req := <-c.SendingCh:
			c.port.Write(req.frame[:])
			select {
			case f := <-c.recvCh:
				req.replyCh <- f
			case <-time.After(c.timeout):
				util.IotLogErrorWithFormatStr("Controller: wait reply timeout for frame: % X", req.frame)
			}
			close(req.replyCh)
		}
	}
}

// RecvCh 暴露板子应答通道，主循环在 HandleSerialNodeMsg 更新状态后向此转发。
func (c *Controller) RecvCh() chan<- serial.Frame {
	return c.recvCh
}

func (c *Controller) sendAndWait(frame [serial.FrameLen]byte) (serial.Frame, bool) {
	replyCh := make(chan serial.Frame, 1)
	req := sendReq{frame: frame, replyCh: replyCh}
	select {
	case c.SendingCh <- req:
	default:
		util.IotLogErrorStr("Controller: SendingCh full, drop request")
		return serial.Frame{}, false
	}
	f, ok := <-replyCh
	return f, ok
}

// ChangeColorForZones 发送指令 1，携带完整 16 区颜色负载。
func (c *Controller) ChangeColorForZones(zones [16]byte) (serial.Frame, bool) {
	payload := serial.PayloadFromZones(zones)
	frame := serial.BuildFrame(serial.CmdChangeColor, payload)
	return c.sendAndWait(frame)
}

// QueryStatus 发送指令 2，返回 16 区状态应答。
func (c *Controller) QueryStatus() (serial.Frame, bool) {
	var payload [serial.PayloadLen]byte
	frame := serial.BuildFrame(serial.CmdQueryStatus, payload)
	return c.sendAndWait(frame)
}

// UpdateNodeStatesFromSerialReply 把板子应答(cmd=1 改色响应 / cmd=2 状态回报)
// 应用到节点状态。应答负载打包了 16 区状态(A/B/C/F)，这里按每个 node 的
// 区范围推导 CompletionStatus。注意：应答里是"状态"而不是"颜色"，颜色仍由
// 网关维护(上次下发的 requesting color 即 reported color，完成态时一致)。
func (c *Controller) UpdateNodeStatesFromSerialReply(frame serial.Frame) {
	zones := serial.ZonesFromPayload(frame.Payload)
	apply := func(nodeList []string) {
		for _, nodeId := range nodeList {
			start, end, ok := NodeToZoneRange(nodeId)
			if !ok {
				continue
			}
			ns := bsp.GetOrCreateNodeState(nodeId)
			if ns == nil {
				continue
			}
			ns.LastMsgTimestamp = time.Now().Unix()
			if AllZonesCompleted(zones, start, end) {
				ns.CompletionStatus = 2
			} else {
				inProgress := false
				for i := start; i < end; i++ {
					if zones[i] == serial.ZoneExecuting || zones[i] == serial.ZonePowering {
						inProgress = true
						break
					}
				}
				if inProgress {
					ns.CompletionStatus = 1
				} else {
					ns.CompletionStatus = 0
				}
			}
		}
	}
	apply(bsp.BspConfigInstance.BaseConfigParams.NodeList1)
	apply(bsp.BspConfigInstance.BaseConfigParams.NodeList2)
}
