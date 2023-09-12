package bsp

import (
	"encoding/hex"
	"encoding/json"
	"iot_go/pkg/shared"
	"iot_go/pkg/thingsboard_shared"
	"iot_go/pkg/util"
	"os"

	"github.com/spf13/viper"
)

type NodeState struct {
	NodeId            string `json:"node_id"`
	LastMsgTimestamp  int64  `json:"node_msg_timestamp"`
	NodeReportedColor []int  `json:"node_reported_color"`
	// NodeRequestingColor []int  `json:"node_requesting_color"`
}

type BspConfig struct {
	shared.InitMsgContent
	shared.BaseConfigParams
	shared.MqttParams
	thingsboard_shared.DeviceProfile `json:"device_profile"`
	NodeStates                       []NodeState `json:"node_state_list"`
}

var swVersion = "0.1.0"

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
		GatewayNodeId: "000000F23456",
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

func DecodeId(i string) []byte {
	b, err := hex.DecodeString(i)
	if err != nil {
		util.IotLogError(err)
	}
	return b
}

func InitConfig() {
	BspConfigInstance.MqttParams.MqttIP = "115.159.53.168"
	BspConfigInstance.MqttParams.MqttPort = 1883
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
