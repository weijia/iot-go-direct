package thingsboard_shared

type ThingsboardServerInfo struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
}

type DeviceProfile struct {
	ThingsboardServerInfo
	ProvisionKey    string `json:"provision_key"`
	ProvisionSecret string `json:"provision_secret"`
}
