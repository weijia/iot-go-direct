package serial_gateway

import (
	"strings"
	"time"

	"go.bug.st/serial"
	"iot_go/pkg/mqtt_util"
)

// Config 网关配置，JSON 字段与 mqtt_serial_gateway-prod.json 对应。
// 内嵌 mqtt_util.MqttParams，因此 mqtt_ip / mqtt_port / mqtt_user_name / mqtt_pwd
// 可直接写在同一份 JSON 里。
type Config struct {
	mqtt_util.MqttParams
	SerialBackend  string `json:"serial_backend"` // real(默认) / echo / generate / simulate，调试用
	SerialPort     string `json:"serial_port"`
	BaudRate       int    `json:"baud_rate"`
	Parity         string `json:"parity"`
	DataBits       int    `json:"data_bits"`
	StopBits       int    `json:"stop_bits"`
	ReadTimeoutMs  int    `json:"read_timeout_ms"`
	FlushMs        int    `json:"flush_ms"`
	ToSerialTopic  string `json:"to_serial_topic"`
	FromSerialTopic string `json:"from_serial_topic"`
	ClientID       string `json:"client_id"`
}

func (c *Config) parity() serial.Parity {
	switch strings.ToUpper(strings.TrimSpace(c.Parity)) {
	case "O":
		return serial.OddParity
	case "E":
		return serial.EvenParity
	case "M":
		return serial.MarkParity
	case "S":
		return serial.SpaceParity
	default:
		return serial.NoParity
	}
}

func (c *Config) stopBits() serial.StopBits {
	switch c.StopBits {
	case 2:
		return serial.TwoStopBits
	case 15, 0:
		return serial.OnePointFiveStopBits
	default:
		return serial.OneStopBit
	}
}

func (c *Config) readTimeout() time.Duration {
	if c.ReadTimeoutMs <= 0 {
		return 50 * time.Millisecond
	}
	return time.Duration(c.ReadTimeoutMs) * time.Millisecond
}

func (c *Config) flushInterval() time.Duration {
	if c.FlushMs <= 0 {
		return 50 * time.Millisecond
	}
	return time.Duration(c.FlushMs) * time.Millisecond
}
