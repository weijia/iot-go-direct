package mqtt_util

type MqttParams struct {
	MqttIP       string `json:"mqtt_ip"`
	MqttPort     int    `json:"mqtt_port"`
	MqttUserName string `json:"mqtt_user_name"`
	MqttPwd      string `json:"mqtt_pwd"`
}
