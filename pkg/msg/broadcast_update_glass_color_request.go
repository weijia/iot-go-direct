package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
	"iot_go/pkg/util"
)

type BroadcastUpdateGlassColorRequest struct {
	Method string      `json:"method"`
	Params ColorParams `json:"params"`
}
type ColorParams struct {
	ColorParams string `json:"color"`
}

func (request BroadcastUpdateGlassColorRequest) handle() interface{} {
	if request.Params.ColorParams != "" {
		// 广播：所有 16 区设置为同一颜色
		color := controller.ParseColorNibble(request.Params.ColorParams)
		var zones [16]byte
		for i := range zones {
			zones[i] = color
		}
		controller.Ctrl.ChangeColorForZones(zones)
		bsp.SetAllRequestingColor(request.Params.ColorParams)
		if frame, ok := controller.Ctrl.QueryStatus(); ok {
			controller.Ctrl.UpdateNodeStatesFromSerialReply(frame)
		}
	} else {
		util.IotLogErrorWithFormatStr("ColorParam is invalid: %v", request.Params.ColorParams)
	}
	reply := GatewayNodeIdReply{
		MsgType:       "broadcast_update_glass_color_reply",
		GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
	}
	return reply
}
