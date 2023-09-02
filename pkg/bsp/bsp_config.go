package bsp

import (
	"encoding/json"
	"iot_go/pkg/shared"
	"iot_go/pkg/thingsboard_shared"
	"iot_go/pkg/util"
	"os"

	"github.com/spf13/viper"
)

type BspConfig struct {
	shared.InitMsgContent
	shared.ConfigParams
	shared.MqttParams
	thingsboard_shared.DeviceProfile `json:"device_profile"`
}

var swVersion = "0.1.0"

var module0Param = shared.Module{
	Freq:   443,
	Band:   250,
	Factor: 10,
}
var module1Param = shared.Module{
	Freq:   444,
	Band:   250,
	Factor: 10,
}
var module2Param = shared.Module{
	Freq:   446,
	Band:   250,
	Factor: 10,
}
var defaultInitMsgContent = shared.InitMsgContent{
	NodeInfoContent: shared.NodeInfoContent{
		GatewayNodeID: "testing-20230815",
		HardVersion:   "0.1.0",
		SoftVersion:   swVersion,
		Custom:        "test",
		Project:       "test",
		NodeType:      1,
		Rssi:          10,
		Ccid:          "test",
		HeartBeat:     20,
		Module1:       module1Param,
		Module2:       module2Param,
	},
	Module0: module0Param,
}
var BspConfigInstance BspConfig

const CONFIG_FILE_NAME = "iot_go.json"

func InitConfig() {
	BspConfigInstance.MqttParams.MqttIP = "115.159.53.168"
	BspConfigInstance.MqttParams.MqttPort = 1883
	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		BspConfigInstance.InitMsgContent = defaultInitMsgContent
		defaultLocalConfig, _ := json.Marshal(BspConfigInstance)
		/*******************  使用 ioutil.WriteFile 写入文件 *****************/
		err2 := os.WriteFile(CONFIG_FILE_NAME, defaultLocalConfig, 0666) //写入文件(字节数组)
		if err2 != nil {
			util.IotLogFatal(err2)
		}

		secondErr := viper.ReadInConfig()
		if secondErr != nil {
			util.IotLogFatal(secondErr)
		}
		viper.WriteConfig()
	}
	data, err := os.ReadFile(CONFIG_FILE_NAME)
	if err == nil && data != nil {
		err = json.Unmarshal(data, &BspConfigInstance)
		if err != nil {
			util.IotLogFatal(err)
		}
	}
}

func (bspConfig BspConfig) CommitChanges() {
	defaultLocalConfig, _ := json.MarshalIndent(bspConfig, "", "    ")
	/*******************  使用 ioutil.WriteFile 写入文件 *****************/
	err2 := os.WriteFile(CONFIG_FILE_NAME, defaultLocalConfig, 0666) //写入文件(字节数组)
	if err2 != nil {
		util.IotLogFatal(err2)
	}
}
