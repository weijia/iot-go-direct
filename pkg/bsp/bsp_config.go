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
	shared.BaseConfigParams
	shared.MqttParams
	thingsboard_shared.DeviceProfile `json:"device_profile"`
	// The following will be real node state, it may contain nodes that is not sent from server
	NodeStates []NodeState `json:"node_state_list"`
	util.LogConfigParams
}

var swVersion = "1.0"

var module0Param = shared.Module{
	Freq:   4723,
	Band:   250,
	Factor: 9,
}
var module1Param = shared.Module{
	Freq:   4723,
	Band:   250,
	Factor: 9,
}
var module2Param = shared.Module{
	Freq:   4463,
	Band:   250,
	Factor: 9,
}

var defaultInitMsgContent = shared.InitMsgContent{
	NodeInfoContent: shared.NodeInfoContent{
		GatewayNodeId: "F12309150001",
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
	// TODO: Change default value to 120.xx.x.x etc
	BspConfigInstance.MqttParams.MqttIP = "app.kosglass.com"
	// BspConfigInstance.MqttParams.MqttIP = "115.159.53.168"
	BspConfigInstance.MqttParams.MqttPort = 1883
	BspConfigInstance.MqttParams.MqttUserName = "l8juew73i2t17wavzthg"
	BspConfigInstance.MqttParams.MqttPwd = "i0eprmhypu3r16g3wuuc"

	BspConfigInstance.DeviceProfile.ProvisionKey = "0hsh1hpc605g4kwyal46"
	BspConfigInstance.DeviceProfile.ProvisionSecret = "68rsgqafhw0anhcwnccr"

	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		BspConfigInstance.InitMsgContent = defaultInitMsgContent
		// BspConfigInstance.ConfigParams.Module1 = module1Param
		// BspConfigInstance.ConfigParams.Module2 = module2Param
		defaultLocalConfig, _ := json.Marshal(BspConfigInstance)
		/*******************  使用 ioutil.WriteFile 写入文件 *****************/
		err2 := os.WriteFile(CONFIG_FILE_NAME, defaultLocalConfig, 0666) //写入文件(字节数组)
		if err2 != nil {
			util.IotLogError(err2)
		}

		secondErr := viper.ReadInConfig()
		if secondErr != nil {
			util.IotLogError(secondErr)
		}
		viper.WriteConfig()
	}
	data, err := os.ReadFile(CONFIG_FILE_NAME)
	if err == nil && data != nil {
		err = json.Unmarshal(data, &BspConfigInstance)
		if err != nil {
			util.IotLogError(err)
		}
	}
}

func (bspConfig BspConfig) CommitChanges() {
	defaultLocalConfig, _ := json.MarshalIndent(bspConfig, "", "    ")
	/*******************  使用 ioutil.WriteFile 写入文件 *****************/
	err2 := os.WriteFile(CONFIG_FILE_NAME, defaultLocalConfig, 0666) //写入文件(字节数组)
	if err2 != nil {
		util.IotLogError(err2)
	}
}
