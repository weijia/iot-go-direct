package util

type IotConfig interface {
	getBroker() string
	getClientId() string
}

type FixedIotConfig struct {
	Broker   string
	ClientId string
}

func (fixedIotConfig FixedIotConfig) GetBroker() string {
	return fixedIotConfig.Broker
}

func (fixedIotConfig FixedIotConfig) GetClientId() string {
	return fixedIotConfig.ClientId
}
