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
	defer ticker.Stop()
	ledFlickerTicker := time.NewTicker(time.Second * time.Duration(1))
	defer ledFlickerTicker.Stop()
	isRefreshed := false
	isLedOn := false
	TurnOnLed(ledManager.DeviceName)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			util.IotLogInfo("Exit led manager due to context done")
			return
		case <-*ledManager.HeartbeatCh:
			// util.IotLogInfo("Received heartbeat for led")
			isRefreshed = true

		case <-ticker.C:
			// TODO: fix led/watch dog calculation issue, it is not 60 second false to trigger
			// TODO: the error print below
			// util.IotLogInfo("led timeout")
			if isRefreshed {
				isRefreshed = false
			}
		case <-ledFlickerTicker.C:
			if isRefreshed {
				if isLedOn {
					TurnOffLed(ledManager.DeviceName)
					isLedOn = false
				} else {
					TurnOnLed(ledManager.DeviceName)
					isLedOn = true
				}
			} else {
				util.IotLogErrorStr("Led manager refresh stopped, led will not flicker anymore")
			}
		}
	}
}
