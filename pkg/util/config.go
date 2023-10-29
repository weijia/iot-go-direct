package util

type IotConfig interface {
	getBroker() string
	getClientId() string
}

type FixedIotConfig struct {
	Broker        string
	SwVersion     string
	GatewayNodeId string
	HardVersion   string
	Ccid          string
	Heartbeat     int
}

func (fixedIotConfig FixedIotConfig) GetBroker() string {
	return fixedIotConfig.Broker
}

func (fixedIotConfig FixedIotConfig) GetClientId() string {
	return fixedIotConfig.GatewayNodeId
}

func (fixedIotConfig FixedIotConfig) GetGatewayNodeId() string {
	return fixedIotConfig.GatewayNodeId
}

func (fixedIotConfig FixedIotConfig) GetHardwareVersion() string {
	return fixedIotConfig.HardVersion
}

func (fixedIotConfig FixedIotConfig) GetCcid() string {
	return fixedIotConfig.Ccid
}

func (fixedIotConfig FixedIotConfig) GetHeartbeat() int {
	return fixedIotConfig.Heartbeat
}
