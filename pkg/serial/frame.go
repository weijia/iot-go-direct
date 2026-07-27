package serial

// 11 字节定长串口帧，协议见 SERIAL_PROTOCOL.md
//   byte0      : 帧头 0x5A
//   byte1      : (指令/状态码)<<4 | 负载长度(固定 8)
//   byte2..9   : 8 字节负载，每字节高低 4 位各控 1 个分区，共 16 区
//   byte10     : 校验 = sum(byte0..9) & 0xFF

const (
	FrameHeader = 0x5A
	FrameLen    = 11
	PayloadLen  = 8

	// 下行指令码（网关 -> 通讯板）
	CmdChangeColor = 1
	CmdQueryStatus = 2

	// 上行状态码（通讯板 -> 网关）
	StatusError       = 0 // 指令错误，负载全 0xFF
	StatusChangeColor = 1 // 改色指令正常响应
	StatusQueryStatus = 2 // 状态回报指令正常响应

	// 下行颜色 nibble 取值
	ColorTransparent = 0x0 // 透明
	ColorLight       = 0x2 // 浅色
	ColorDark        = 0x4 // 深色
	ColorKeep        = 0xF // 保留原有颜色，不变更

	// 上行分区状态 nibble 取值
	ZoneCompleted  = 0xA // 已完成
	ZoneExecuting  = 0xB // 执行中
	ZonePowering   = 0xC // 上电中
	ZoneStatusKeep = 0xF // 保留不变
)

type Frame struct {
	Cmd     byte // byte1 高 4 位：指令码或状态码
	Payload [PayloadLen]byte
	Raw     [FrameLen]byte
}

func checksum(b []byte) byte {
	var sum byte
	for _, x := range b {
		sum += x
	}
	return sum & 0xFF
}

// BuildFrame assembles an 11-byte frame from a command code and 8-byte payload.
func BuildFrame(cmd byte, payload [PayloadLen]byte) [FrameLen]byte {
	var f [FrameLen]byte
	f[0] = FrameHeader
	f[1] = (cmd << 4) | PayloadLen
	for i := 0; i < PayloadLen; i++ {
		f[2+i] = payload[i]
	}
	f[FrameLen-1] = checksum(f[:FrameLen-1])
	return f
}

// ParseFrame validates and decodes a raw 11-byte frame.
func ParseFrame(raw []byte) (Frame, bool) {
	if len(raw) != FrameLen {
		return Frame{}, false
	}
	if raw[0] != FrameHeader {
		return Frame{}, false
	}
	if checksum(raw[:FrameLen-1]) != raw[FrameLen-1] {
		return Frame{}, false
	}
	var f Frame
	f.Cmd = raw[1] >> 4
	copy(f.Payload[:], raw[2:2+PayloadLen])
	copy(f.Raw[:], raw)
	return f, true
}

// PayloadFromZones packs 16 个分区 nibble 到 8 字节负载。
// zones[0]=1 区(字节0高 4 位), zones[1]=2 区(字节0低 4 位), ...
func PayloadFromZones(zones [16]byte) [PayloadLen]byte {
	var p [PayloadLen]byte
	for i := 0; i < PayloadLen; i++ {
		p[i] = (zones[2*i] << 4) | (zones[2*i+1] & 0xF)
	}
	return p
}

// ZonesFromPayload 是 PayloadFromZones 的逆操作。
func ZonesFromPayload(p [PayloadLen]byte) [16]byte {
	var z [16]byte
	for i := 0; i < PayloadLen; i++ {
		z[2*i] = (p[i] >> 4) & 0xF
		z[2*i+1] = p[i] & 0xF
	}
	return z
}
