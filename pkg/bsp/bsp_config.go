package bsp

import (
	"bufio"
	"encoding/json"
	"iot_go/pkg/shared"
	"iot_go/pkg/thingsboard_shared"
	"iot_go/pkg/util"
	"os"
	"path/filepath"
	"strings"

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

var SwVersion = util.APP_VERSION

var module0Param = shared.Module{
	Freq:   4723,
	Band:   250,
	Factor: 9,
}
var module1Param = shared.Module{
	Freq:   4734,
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
		HardVersion:   "1.0",
		SoftVersion:   SwVersion,
		Custom:        "test",
		Project:       "test",
		NodeType:      1,
		Rssi:          10,
		Ccid:          "test",
		Heartbeat:     20,
		Module1:       module1Param,
		Module2:       module2Param,
	},
	Module0: module0Param,
}
var BspConfigInstance BspConfig

const CONFIG_FILE_NAME = "iot_go.json"

func InitConfig() {
	// TODO: Change default value to 120.xx.x.x etc
	// BspConfigInstance.MqttParams.MqttIP = "app.kosglass.com"
	BspConfigInstance.MqttParams.MqttIP = "120.79.55.61"
	// BspConfigInstance.MqttParams.MqttIP = "115.159.53.168"
	BspConfigInstance.MqttParams.MqttPort = 1883
	BspConfigInstance.MqttParams.MqttUserName = "l8juew73i2t17wavzthg"
	BspConfigInstance.MqttParams.MqttPwd = "i0eprmhypu3r16g3wuuc"

	BspConfigInstance.DeviceProfile.ProvisionKey = "0hsh1hpc605g4kwyal46"
	BspConfigInstance.DeviceProfile.ProvisionSecret = "68rsgqafhw0anhcwnccr"

	viper.SetConfigName(CONFIG_FILE_NAME)
	viper.SetConfigType("json")
	viper.AddConfigPath(util.GetAppRoot())
	err := viper.ReadInConfig()
	appRoot := util.GetAppRoot()
	configFilePath := filepath.Join(appRoot, CONFIG_FILE_NAME)

	if err != nil {
		util.IotLogErrWithStr("vip read config error", err)
		BspConfigInstance.InitMsgContent = defaultInitMsgContent
		// BspConfigInstance.ConfigParams.Module1 = module1Param
		// BspConfigInstance.ConfigParams.Module2 = module2Param
		defaultLocalConfig, _ := json.Marshal(BspConfigInstance)
		/*******************  使用 ioutil.WriteFile 写入文件 *****************/
		err2 := os.WriteFile(configFilePath, defaultLocalConfig, 0666) //写入文件(字节数组)
		if err2 != nil {
			util.IotLogError(err2)
		}

		secondErr := viper.ReadInConfig()
		if secondErr != nil {
			util.IotLogError(secondErr)
		}
		viper.WriteConfig()
	}
	data, err := os.ReadFile(configFilePath)
	if err == nil && data != nil {
		err = json.Unmarshal(data, &BspConfigInstance)
		if err != nil {
			util.IotLogError(err)
		}
	}
	// Overwrite board related info to config
	BspConfigInstance.SoftVersion = SwVersion
	gatewayIdFilePath := filepath.Join(appRoot, "gateway_id.txt")
	file, err := os.Open(gatewayIdFilePath)

	if err == nil {
		defer file.Close()
		r := bufio.NewReader(file)
		line, _, e := r.ReadLine()
		if e == nil {
			s := strings.Replace(string(line), "\n", "", -1)
			s = strings.Replace(s, "\r", "", -1)
			s = strings.Replace(s, " ", "", -1)
			util.IotLog("Len: %d", len(s))
			if len(s) == 12 && len(util.DecodeId(s)) == 6 {
				BspConfigInstance.GatewayNodeId = s
			} else {
				util.IotLogErrorStr("gateway_id.txt contain invalid gateway id: " + string(data))
			}
		}
	} else {
		util.IotLogErrWithStr("Open gateway_id.txt failed", err)
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
