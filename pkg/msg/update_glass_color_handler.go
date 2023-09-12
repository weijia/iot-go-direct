package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

type UpdateGlassColorRequest struct {
	Method string                          `json:"method"`
	Params []shared.UpdateGlassColorParams `json:"params"`
}

func (request UpdateGlassColorRequest) handle(mqttClient *util.Mqtt) interface{} {
	reply := bsp.GetBsp().SetGlassColors(mqttClient, request.Params)

	return nil
}
