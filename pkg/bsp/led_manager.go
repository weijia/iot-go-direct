package bsp

import (
	"context"
	"iot_go/pkg/util"
	"time"
)

type LedManager struct {
	HeartbeatCh *chan int
	DeviceName  string
	Timeout     int
}

func (ledManager LedManager) LedMsgLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(ledManager.Timeout))
	isRefreshed := false
	TurnOnLed(ledManager.DeviceName)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			util.IotLogInfo("Exit led manager due to contex done")
			return
		case <-*ledManager.HeartbeatCh:
			// util.IotLogInfo("Received heartbeat for led")
			isRefreshed = true

		case <-ticker.C:
			// util.IotLogInfo("led timeout")
			if !isRefreshed {
				util.IotLogErrorStr("LED refresh timeout")
				TurnOffLed(ledManager.DeviceName)
			} else {
				TurnOnLed(ledManager.DeviceName)
				isRefreshed = false
			}
		}
	}
}
