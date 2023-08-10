package util

type IotConfig interface {
	getBroker() string
	getClientId() string
}

type FixedIotConfig struct {
	Broker        string
	SwVersion     string
	GatewayNodeID string
	HardVersion   string
	Ccid          string
	HeartBeat     int
}

func (fixedIotConfig FixedIotConfig) GetBroker() string {
	return fixedIotConfig.Broker
}

func (fixedIotConfig FixedIotConfig) GetClientId() string {
	return fixedIotConfig.GatewayNodeID
}

func (fixedIotConfig FixedIotConfig) GetGatewayNodeID() string {
	return fixedIotConfig.GatewayNodeID
}

func (fixedIotConfig FixedIotConfig) GetHardwareVersion() string {
	return fixedIotConfig.HardVersion
}

func (fixedIotConfig FixedIotConfig) GetCcid() string {
	return fixedIotConfig.Ccid
}

func (fixedIotConfig FixedIotConfig) GetHeartBeat() int {
	return fixedIotConfig.HeartBeat
}
