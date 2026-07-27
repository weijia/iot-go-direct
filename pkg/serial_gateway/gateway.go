package serial_gateway

import (
	"context"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"iot_go/pkg/mqtt_util"
	"iot_go/pkg/util"
)

const (
	recvChannelSize = 256
	serialReadBuf   = 1024
	upstreamMaxBuf  = 4096
)

// Gateway 桥接 MQTT 与串口，采用原始字节直通（raw passthrough）：
//   - 下行：订阅 to_serial_topic，收到消息直接写入串口；
//   - 上行：从串口读字节，按 flush 间隔或缓冲满批量发布到 from_serial_topic。
type Gateway struct {
	cfg    Config
	mqtt   *mqtt_util.MqttEasyClient
	serial SerialConn
}

// NewGateway 创建并打开串口、构造 MQTT 客户端（此时尚未建立连接）。
func NewGateway(cfg Config) (*Gateway, error) {
	recvCh := make(chan mqtt.Message, recvChannelSize)
	easyClient := &mqtt_util.MqttEasyClient{
		MqttParams:       cfg.MqttParams,
		MqttClientId:     cfg.ClientID,
		Topic:            cfg.ToSerialTopic,
		ReceivingChannel: &recvCh,
	}
	sp, err := OpenSerialConn(cfg)
	if err != nil {
		return nil, err
	}
	return &Gateway{cfg: cfg, mqtt: easyClient, serial: sp}, nil
}

// Run 启动桥接：后台连接 MQTT（自带自动重连与重订阅），并起上下行两个协程。
func (g *Gateway) Run(ctx context.Context) error {
	go func() {
		if err := g.mqtt.ConnectAndSubscribe(); err != nil {
			util.IotLogErrWithStr("mqtt connect", err)
		}
	}()
	go g.downstreamLoop(ctx)
	go g.upstreamLoop(ctx)
	return nil
}

// downstreamLoop：MQTT -> 串口（原始字节直写）
func (g *Gateway) downstreamLoop(ctx context.Context) {
	ch := g.mqtt.ReceivingChannel
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-*ch:
			if _, err := g.serial.Write(msg.Payload()); err != nil {
				util.IotLogErrWithStr("serial write", err)
			} else {
				util.IotLog(">> serial: %d bytes", len(msg.Payload()))
			}
		}
	}
}

// upstreamLoop：串口 -> MQTT（按 flush 间隔或缓冲满批量发布）
func (g *Gateway) upstreamLoop(ctx context.Context) {
	buf := make([]byte, 0, upstreamMaxBuf)
	readBuf := make([]byte, serialReadBuf)
	ticker := time.NewTicker(g.cfg.flushInterval())
	defer ticker.Stop()

	flush := func() {
		if len(buf) == 0 {
			return
		}
		g.publish(buf)
		buf = buf[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case <-ticker.C:
			flush()
		default:
		}
		n, err := g.serial.Read(readBuf)
		if err != nil {
			util.IotLogErrWithStr("serial read", err)
			flush()
			// 尝试重新打开真实串口（仅真实后端支持），失败则退避后继续
			if sp, ok := g.serial.(*SerialPort); ok {
				if rerr := sp.reopen(); rerr != nil {
					time.Sleep(time.Second)
				}
			} else {
				time.Sleep(time.Second)
			}
			continue
		}
		if n > 0 {
			buf = append(buf, readBuf[:n]...)
			if len(buf) >= upstreamMaxBuf {
				flush()
			}
		}
	}
}

func (g *Gateway) publish(data []byte) {
	if g.mqtt.Client == nil {
		return
	}
	token := g.mqtt.Client.Publish(g.cfg.FromSerialTopic, 0, false, data)
	token.Wait()
	if token.Error() != nil {
		util.IotLogErrWithStr("mqtt publish", token.Error())
	} else {
		util.IotLog("<< mqtt: %d bytes to %s", len(data), g.cfg.FromSerialTopic)
	}
}

// Close 优雅关闭串口与 MQTT 连接。
func (g *Gateway) Close() {
	if g.serial != nil {
		_ = g.serial.Close()
	}
	if g.mqtt != nil && g.mqtt.Client != nil {
		g.mqtt.Client.Disconnect(250)
	}
}
