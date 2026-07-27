package msg

import (
	"iot_go/pkg/bsp"
	"iot_go/pkg/controller"
	"iot_go/pkg/util"
)

type GroupUpdateGlassColorRequest struct {
	Method string                      `json:"method"`
	Params GroupUpdateGlassColorParams `json:"params"`
}
type GroupUpdateGlassColorParams struct {
	NodeList1 []string `json:"node_list1"`
	NodeList2 []string `json:"node_list2"`
	Color     string   `json:"color"`
}

func (request GroupUpdateGlassColorRequest) Replay() {
	if groupUpdateGlassColorRequestResendCnt > 0 {
		util.IotLogErrorWithFormatStr("resend group update glass color request, cnt: %d", groupUpdateGlassColorRequestResendCnt)
		groupUpdateGlassColorRequestResendCnt -= 1
		request.handle()
	}
}

func (request GroupUpdateGlassColorRequest) handle() interface{} {
	invalidNodeList := []string{}
	group := request.Params

	color := controller.ParseColorNibble(group.Color)
	var zones [16]byte
	for i := range zones {
		zones[i] = controller.ParseColorNibble("f") // 默认保留(F)，仅覆盖列出的节点
	}

	for _, nodeId := range group.NodeList1 {
		if bsp.IsInNodeList1(nodeId) {
			bsp.SetRequestingColor(nodeId, group.Color)
			for i := 0; i < 8; i++ {
				zones[i] = color
			}
		} else {
			invalidNodeList = append(invalidNodeList, nodeId)
		}
	}
	for _, nodeId := range group.NodeList2 {
		if bsp.IsInNodeList2(nodeId) {
			bsp.SetRequestingColor(nodeId, group.Color)
			for i := 8; i < 16; i++ {
				zones[i] = color
			}
		} else {
			invalidNodeList = append(invalidNodeList, nodeId)
		}
	}

	controller.Ctrl.ChangeColorForZones(zones)
	if frame, ok := controller.Ctrl.QueryStatus(); ok {
		controller.Ctrl.UpdateNodeStatesFromSerialReply(frame)
	}

	var reply GroupUpdateGlassColorReply
	reply.MsgType = "group_update_glass_color_reply"
	reply.GatewayNodeId = bsp.BspConfigInstance.GatewayNodeId
	reply.InvalidNodes = invalidNodeList
	return reply
}
