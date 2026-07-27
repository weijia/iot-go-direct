package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
	"iot_go/pkg/node"
	"iot_go/pkg/serial"
	"iot_go/pkg/shared"
	"sync"
)

type UpdateGlassColorRequest struct {
	Method string                          `json:"method"`
	Params []shared.UpdateGlassColorParams `json:"params"`
}

func (request UpdateGlassColorRequest) handle(mqttToServer chan interface{}) {
	var reply shared.UpdateGlassColorReply

	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.MsgType = "update_glass_color_reply"
	reply.Status = make([]shared.UpdateGlassColorParams, 0, len(request.Params))

	var wg sync.WaitGroup
	for _, param := range request.Params {
		wg.Add(1)
		go func(param shared.UpdateGlassColorParams) {
			defer wg.Done()
			// 记录请求颜色（颜色由网关自己维护，完成态时即上报颜色）
			bsp.SetRequestingColor(param.NodeId, param.Color)
			zones, ok := controller.NodeColorToZones(param.NodeId, param.Color)
			if !ok {
				reply.Status = append(reply.Status, shared.UpdateGlassColorParams{
					NodeId: param.NodeId,
					Color:  node.SetColorForNodeAsInvalid(param.Color),
				})
				return
			}
			frame, sent := controller.Ctrl.ChangeColorForZones(zones)
			if !sent || frame.Cmd != serial.StatusChangeColor {
				// 板子无应答或回包异常（如上电中阻塞导致超时）
				reply.Status = append(reply.Status, shared.UpdateGlassColorParams{
					NodeId: param.NodeId,
					Color:  node.SetColorForNodeAsInvalid(param.Color),
				})
				return
			}
			controller.Ctrl.UpdateNodeStatesFromSerialReply(frame)
			reply.Status = append(reply.Status, param)
		}(param)
	}
	wg.Wait()

	mqttToServer <- reply
}
