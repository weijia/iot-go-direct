package serial_gateway

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"

	"iot_go/pkg/util"
)

// SerialConn 串口连接抽象。真实串口和调试用的假串口都实现它，
// 这样网关核心逻辑不依赖具体硬件，便于无设备调试与单元测试。
type SerialConn interface {
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Close() error
}

func openPort(config Config) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: config.BaudRate,
		DataBits: config.DataBits,
		StopBits: config.stopBits(),
		Parity:   config.parity(),
	}
	port, err := serial.Open(config.SerialPort, mode)
	if err != nil {
		return nil, err
	}
	if err := port.SetReadTimeout(config.readTimeout()); err != nil {
		util.IotLogErrWithStr("set serial read timeout", err)
	}
	return port, nil
}

// SerialPort 真实串口封装（基于 go.bug.st/serial）。
type SerialPort struct {
	port   serial.Port
	config Config
}

func OpenSerial(config Config) (*SerialPort, error) {
	port, err := openPort(config)
	if err != nil {
		return nil, err
	}
	return &SerialPort{port: port, config: config}, nil
}

func (s *SerialPort) Read(p []byte) (int, error)  { return s.port.Read(p) }
func (s *SerialPort) Write(p []byte) (int, error) { return s.port.Write(p) }
func (s *SerialPort) Close() error                { return s.port.Close() }

// reopen 在串口异常后尝试关闭并重新打开（保持原有配置）。
func (s *SerialPort) reopen() error {
	_ = s.port.Close()
	port, err := openPort(s.config)
	if err != nil {
		return err
	}
	s.port = port
	return nil
}

// OpenSerialConn 按配置选择串口后端：
//   - "echo"      调试用：把写入的字节回环到读取端（验证下行→上行闭环）
//   - "generate"  调试用：定时产生测试字节（验证上行→MQTT 路径）
//   - "simulate"  调试用：模拟被控玻璃通讯板，按 SERIAL_PROTOCOL 校验下行包并构造合法应答回发
//   - 其它/空     使用真实串口（serial_port 字段指定的 COM/tty 设备）
func OpenSerialConn(config Config) (SerialConn, error) {
	switch strings.ToLower(config.SerialBackend) {
	case "echo":
		util.IotLog("Using echo fake serial backend (debug)")
		return newEchoSerial(), nil
	case "generate":
		util.IotLog("Using generate fake serial backend (debug)")
		interval := config.flushInterval()
		if interval <= 0 {
			interval = time.Second
		}
		return newGenerateSerial(interval, []byte("DEBUG")), nil
	case "simulate":
		util.IotLog("Using simulate serial backend (acts as glass board, replies per protocol)")
		return newSimulateSerial(), nil
	default:
		return OpenSerial(config)
	}
}

// echoSerial 调试后端：Write 的内容会被 Read 原样返回（回环）。
type echoSerial struct {
	mu     sync.Mutex
	buf    []byte
	cond   *sync.Cond
	closed bool
}

func newEchoSerial() *echoSerial {
	e := &echoSerial{}
	e.cond = sync.NewCond(&e.mu)
	return e
}

func (e *echoSerial) Write(p []byte) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return 0, io.EOF
	}
	e.buf = append(e.buf, p...)
	e.cond.Broadcast()
	return len(p), nil
}

func (e *echoSerial) Read(p []byte) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for len(e.buf) == 0 && !e.closed {
		e.cond.Wait()
	}
	if e.closed && len(e.buf) == 0 {
		return 0, io.EOF
	}
	n := copy(p, e.buf)
	e.buf = e.buf[n:]
	return n, nil
}

func (e *echoSerial) Close() error {
	e.mu.Lock()
	e.closed = true
	e.cond.Broadcast()
	e.mu.Unlock()
	return nil
}

// generateSerial 调试后端：定时产生测试字节，供验证"串口→MQTT"上行路径。
type generateSerial struct {
	ch       chan []byte
	closeCh  chan struct{}
	counter  int
	payload  []byte
	interval time.Duration
}

func newGenerateSerial(interval time.Duration, payload []byte) *generateSerial {
	g := &generateSerial{
		ch:       make(chan []byte, 16),
		closeCh:  make(chan struct{}),
		payload:  payload,
		interval: interval,
	}
	go g.run()
	return g
}

func (g *generateSerial) run() {
	ticker := time.NewTicker(g.interval)
	defer ticker.Stop()
	for {
		select {
		case <-g.closeCh:
			return
		case <-ticker.C:
			g.counter++
			msg := []byte(fmt.Sprintf("%s#%d\n", string(g.payload), g.counter))
			select {
			case g.ch <- msg:
			default:
			}
		}
	}
}

func (g *generateSerial) Read(p []byte) (int, error) {
	select {
	case data := <-g.ch:
		return copy(p, data), nil
	case <-g.closeCh:
		return 0, io.EOF
	}
}

func (g *generateSerial) Write(p []byte) (int, error) { return len(p), nil } // 丢弃
func (g *generateSerial) Close() error {
	close(g.closeCh)
	return nil
}

// simulateSerial 调试后端：模拟「被控玻璃通讯板」。
// 网关下发的字节(Write)被视为 网关→通讯板 的下行指令，这里按 SERIAL_PROTOCOL
// 校验其正确性，并构造合法的 11 字节上行应答包(Read 返回)。从而在不接串口的情况下
// 端到端验证「网页下发 → 网关透传 → 协议校验 → 按协议回复」全链路。
type simulateSerial struct {
	mu     sync.Mutex
	buf    []byte
	cond   *sync.Cond
	closed bool
	zones  [16]byte // 各分区当前状态，默认 A(已完成)
}

func newSimulateSerial() *simulateSerial {
	s := &simulateSerial{}
	for i := range s.zones {
		s.zones[i] = 0x0A // A = 已完成
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *simulateSerial) Write(p []byte) (int, error) {
	resp := s.handleDownlink(p)
	s.mu.Lock()
	if !s.closed && resp != nil {
		s.buf = append(s.buf, resp...)
		s.cond.Broadcast()
	}
	s.mu.Unlock()
	return len(p), nil
}

func (s *simulateSerial) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for len(s.buf) == 0 && !s.closed {
		s.cond.Wait()
	}
	if s.closed && len(s.buf) == 0 {
		return 0, io.EOF
	}
	n := copy(p, s.buf)
	s.buf = s.buf[n:]
	return n, nil
}

func (s *simulateSerial) Close() error {
	s.mu.Lock()
	s.closed = true
	s.cond.Broadcast()
	s.mu.Unlock()
	return nil
}

// handleDownlink 按 SERIAL_PROTOCOL 校验下行包，返回要回的应答包；nil 表示按协议不回复。
func (s *simulateSerial) handleDownlink(p []byte) []byte {
	if len(p) != 11 {
		util.IotLog("simulate: 下行包长度=%d(应11)，按协议丢弃", len(p))
		return nil
	}
	if p[0] != 0x5A {
		util.IotLog("simulate: 帧头=0x%02X(应5A)，按协议丢弃", p[0])
		return nil
	}
	var sum byte
	for i := 0; i < 10; i++ {
		sum += p[i]
	}
	if sum != p[10] {
		util.IotLog("simulate: 校验失败(计算0x%02X≠收到0x%02X)，按协议丢弃", sum, p[10])
		return nil
	}
	cmd := (p[1] >> 4) & 0x0F
	switch cmd {
	case 1, 2:
		// 合法指令：模拟执行后回报。cmd=1 改色 / cmd=2 回报 均回报各分区当前状态。
		util.IotLog("simulate: 收到合法指令 cmd=%d，回复状态码=%d(16区均已完成 A)", cmd, cmd)
		return s.buildResponse(byte(cmd))
	default:
		// 指令码非法：状态码0，负载全 0xFF
		util.IotLog("simulate: 指令码=%d 非法，回复状态码0(负载全FF)", cmd)
		return s.buildErrorResponse()
	}
}

// buildResponse 构造状态码为 statusCode 的应答包，负载上报 16 分区当前状态。
func (s *simulateSerial) buildResponse(statusCode byte) []byte {
	pkt := make([]byte, 11)
	pkt[0] = 0x5A
	pkt[1] = (statusCode << 4) | 0x08 // 低4位长度固定 8
	for i := 0; i < 16; i++ {
		if i%2 == 0 {
			pkt[2+i/2] |= (s.zones[i] & 0x0F) << 4
		} else {
			pkt[2+i/2] |= s.zones[i] & 0x0F
		}
	}
	var sum byte
	for i := 0; i < 10; i++ {
		sum += pkt[i]
	}
	pkt[10] = sum
	return pkt
}

// buildErrorResponse 构造状态码0、负载全 0xFF 的错误应答包。
func (s *simulateSerial) buildErrorResponse() []byte {
	pkt := make([]byte, 11)
	pkt[0] = 0x5A
	pkt[1] = (0 << 4) | 0x08
	for i := 2; i <= 9; i++ {
		pkt[i] = 0xFF
	}
	var sum byte
	for i := 0; i < 10; i++ {
		sum += pkt[i]
	}
	pkt[10] = sum
	return pkt
}
