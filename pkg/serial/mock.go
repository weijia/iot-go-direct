package serial

import (
	"io"
	"log"
	"net"
	"time"
)

const mockExecuteDelay = 800 * time.Millisecond

// MockBoard 在内存管道另一端模拟通讯板，按 SERIAL_PROTOCOL.md 实现：
//   - 收到改色(cmd=1) 把非保留区标记执行中，立即回包(含当前区状态)，异步后置为已完成；
//   - 收到状态回报(cmd=2) 回包当前 16 区状态；
//   - powering 指向 true 时对任何帧静默丢弃、不回包（触发网关超时）。
type MockBoard struct {
	conn     io.ReadWriteCloser
	zoneStat [16]byte
	powering *bool
}

func NewMockBoard(conn io.ReadWriteCloser, powering *bool) *MockBoard {
	b := &MockBoard{conn: conn, powering: powering}
	for i := range b.zoneStat {
		b.zoneStat[i] = ZoneCompleted
	}
	return b
}

func (b *MockBoard) readFrame() ([FrameLen]byte, bool) {
	var buf [FrameLen]byte
	if _, err := io.ReadFull(b.conn, buf[:]); err != nil {
		return buf, false
	}
	return buf, true
}

func (b *MockBoard) reply(status byte) {
	payload := PayloadFromZones(b.zoneStat)
	frame := BuildFrame(status, payload)
	log.Printf("[MockBoard] reply cmd=%d (zones=%v)", status, b.zoneStat)
	b.conn.Write(frame[:])
}

func (b *MockBoard) run() {
	for {
		raw, ok := b.readFrame()
		if !ok {
			return
		}
		f, ok := ParseFrame(raw[:])
		if !ok {
			continue // 非法帧忽略
		}
		if b.powering != nil && *b.powering {
			log.Printf("[MockBoard] powering -> drop frame cmd=%d (no reply)", f.Cmd)
			continue // 上电中：静默丢弃、不回包
		}
		log.Printf("[MockBoard] recv cmd=%d", f.Cmd)
		switch f.Cmd {
		case CmdChangeColor:
			zones := ZonesFromPayload(f.Payload)
			for i := 0; i < 16; i++ {
				if zones[i] != ColorKeep {
					b.zoneStat[i] = ZoneExecuting
				}
			}
			b.reply(StatusChangeColor)
			go func(z [16]byte) {
				time.Sleep(mockExecuteDelay)
				for i := 0; i < 16; i++ {
					if z[i] != ColorKeep {
						b.zoneStat[i] = ZoneCompleted
					}
				}
			}(zones)
		case CmdQueryStatus:
			b.reply(StatusQueryStatus)
		}
	}
}

// OpenMock 创建一对互联的内存传输，并在另一端启动 MockBoard。
// 用于在没有任何串口设备的情况下测试整个串口/控制板通信栈。
//   - powering 为 nil 时板子永远正常应答；
//   - 若传 *bool，可在运行时翻转来模拟"上电中"阻塞(静默丢帧)。
func OpenMock(powering *bool) *Port {
	gw, board := net.Pipe()
	b := NewMockBoard(board, powering)
	go b.run()
	return New(gw, Config{Port: "mock"})
}
