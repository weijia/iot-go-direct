package bsp

import (
	"encoding/json"
	"iot_go/pkg/shared"

	"fmt"
)

type Bsp interface {
	SetModule0Params(module shared.Module)
	SetModule1Params(module shared.Module)
	SetModule2Params(module shared.Module)
}

func InitBoard() {

}

type VirtualBsp struct {
}

func (virtualBsp VirtualBsp) SetModule0Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Println("Setting Module0 Params", moduleParamsInJson)
}

func (virtualBsp VirtualBsp) SetModule1Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Println("Setting Module1 Params", moduleParamsInJson)
}

func (virtualBsp VirtualBsp) SetModule2Params(moduleParams shared.Module) {
	moduleParamsInJson, _ := json.Marshal(moduleParams)
	fmt.Println("Setting Module2 Params", moduleParamsInJson)
}

func GetBsp() VirtualBsp {
	var virtualBsp VirtualBsp
	return virtualBsp
}
