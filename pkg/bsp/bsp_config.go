package bsp

import (
	"encoding/json"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"os"

	"github.com/spf13/viper"
)

type BspConfig struct {
	InitConfig   shared.Init
	Broker       string
	ConfigParams shared.ConfigParams
}

var swVersion = "0.1.0"

var module0 = shared.Module{
	Freq:   443,
	Band:   250,
	Factor: 10,
}
var module1 = shared.Module{
	Freq:   444,
	Band:   250,
	Factor: 10,
}
var module2 = shared.Module{
	Freq:   446,
	Band:   250,
	Factor: 10,
}
var defaultInitMsg = shared.Init{
	MsgType:       "init",
	GatewayNodeID: "test",
	HardVersion:   "0.1.0",
	SoftVersion:   swVersion,
	Custom:        "test",
	Project:       "test",
	NodeType:      1,
	Rssi:          10,
	Ccid:          "test",
	HeartBeat:     20,
	Module0:       module0,
	Module1:       module1,
	Module2:       module2,
}
var BspConfigInstance BspConfig

func InitConfig() {
	BspConfigInstance.Broker = "tcp://115.159.53.168:1883"
	viper.SetConfigName("iot_go.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		BspConfigInstance.InitConfig = defaultInitMsg
		defaultLocalConfig, _ := json.Marshal(BspConfigInstance)
		/*******************  使用 ioutil.WriteFile 写入文件 *****************/
		err2 := os.WriteFile("./iot_go.json", defaultLocalConfig, 0666) //写入文件(字节数组)
		if err2 != nil {
			util.IotLogFatal(err2)
		}

		secondErr := viper.ReadInConfig()
		if secondErr != nil {
			util.IotLogFatal(secondErr)
		}
		viper.WriteConfig()
	}
	data, err := os.ReadFile("./iot_go.json")
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
	err2 := os.WriteFile("./iot_go.json", defaultLocalConfig, 0666) //写入文件(字节数组)
	if err2 != nil {
		util.IotLogFatal(err2)
	}
}
