package serial

import (
	"io"
	"time"

	"go.bug.st/serial"

	"iot_go/pkg/util"
)

type Config struct {
	Port string
	Baud int
}

type Port struct {
	cfg    Config
	conn   io.ReadWriteCloser
	recvCh chan Frame
	quit   chan struct{}
}

// Open opens a real serial port and starts the background read loop.
func Open(cfg Config) (*Port, error) {
	mode := &serial.Mode{
		BaudRate: cfg.Baud,
	}
	if mode.BaudRate == 0 {
		mode.BaudRate = 9600
	}
	conn, err := serial.Open(cfg.Port, mode)
	if err != nil {
		return nil, err
	}
	return New(conn, cfg), nil
}

// New wraps an arbitrary io.ReadWriteCloser as a Port and starts the read loop.
// Pass a real serial.Port for production, or an in-memory pipe (e.g. net.Pipe)
// for testing without hardware.
func New(conn io.ReadWriteCloser, cfg Config) *Port {
	p := &Port{
		cfg:    cfg,
		conn:   conn,
		recvCh: make(chan Frame, 16),
		quit:   make(chan struct{}),
	}
	go p.readLoop()
	return p
}

// readLoop scans the byte stream for 0x5A 起始的 11 字节包，校验通过后推入
// recvCh。使用滑动窗口，避免负载里出现的 0x5A 造成永久失步。
func (p *Port) readLoop() {
	buf := make([]byte, 64)
	var window []byte
	for {
		select {
		case <-p.quit:
			return
		default:
		}
		n, err := p.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			util.IotLogErrorWithFormatStr("serial read error: %v", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}
		for i := 0; i < n; i++ {
			window = append(window, buf[i])
			if len(window) > FrameLen {
				window = window[1:]
			}
			if len(window) == FrameLen {
				if f, ok := ParseFrame(window); ok {
					select {
					case p.recvCh <- f:
					default:
						util.IotLogErrorStr("serial recvCh full, drop frame")
					}
					window = nil
				} else {
					// 非法帧：丢弃最旧字节继续扫描
					window = window[1:]
				}
			}
		}
	}
}

func (p *Port) Write(b []byte) (int, error) {
	return p.conn.Write(b)
}

// RecvCh returns the channel the read loop pushes valid frames to.
func (p *Port) RecvCh() <-chan Frame {
	return p.recvCh
}

func (p *Port) Close() error {
	close(p.quit)
	return p.conn.Close()
}
