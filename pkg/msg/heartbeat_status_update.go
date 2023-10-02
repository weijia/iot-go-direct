package msg

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"time"
)

type HeartbeatStatusUpdate struct {
	MsgType       string            `json:"msg_type"`
	GatewayNodeId string            `json:"gateway_node_id"`
	Status        []HeartbeatStatus `json:"status"`
}
type HeartbeatStatus struct {
	NodeId      string `json:"node_id"`
	Color       string `json:"color"`
	HardVersion string `json:"hard_version"`
	SoftVersion string `json:"soft_version"`
	RunArea     int    `json:"run_area"`
}

func GetColorStrFromSlice(colorSlice [8]int) string {
	res := ""
	for index, color := range colorSlice {
		res += fmt.Sprintf("%x%x", index+1, color)
	}
	return res
}

func SendHeartbeatToServer() {
	// Create a ticker with 50 ms = 0.05 seconds
	ticker := time.NewTicker(time.Second * 60)
	u := HeartbeatStatusUpdate{
		MsgType:       "heatbeat_status_update",
		GatewayNodeId: bsp.BspConfigInstance.GatewayNodeId,
	}
	for {
		<-ticker.C
		currentTimestamp := time.Now().Unix()
		var l []HeartbeatStatus
		for _, state := range bsp.BspConfigInstance.NodeStates {
			if bsp.IsInNodeList1(state.NodeId) || bsp.IsInNodeList2(state.NodeId) {
				color := GetColorStrFromSlice(state.NodeReportedColor)
				if state.LastMsgTimestamp < currentTimestamp-int64(bsp.BspConfigInstance.Heartbeat) {
					color = node.SetColorForNodeAsInvalid(color)
				}
				s := HeartbeatStatus{
					NodeId:      state.NodeId,
					Color:       color,
					HardVersion: fmt.Sprintf("%x", state.HwVer),
					SoftVersion: fmt.Sprintf("%x", state.SwVer),
					RunArea:     state.RunningArea,
				}
				l = append(l, s)
			}
		}
		u.Status = l
		MqttToServerCh <- u
	}
}
