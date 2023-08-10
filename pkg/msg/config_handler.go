package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
)

type ConfigRequest struct {
	Method string              `json:"method"`
	Params shared.ConfigParams `json:"params"`
}

func (config ConfigRequest) handle() {
	bspInstance := bsp.GetBsp()
	bspInstance.SetModule1Params(config.Params.Module1)
	bspInstance.SetModule2Params(config.Params.Module2)
	// fmt.Printf("%s", config.Method)
	bsp.BspConfigInstance.InitConfig.Module1 = config.Params.Module1
	bsp.BspConfigInstance.InitConfig.Module2 = config.Params.Module2
	bsp.BspConfigInstance.ConfigParams = config.Params
	bsp.BspConfigInstance.CommitChanges()
}
