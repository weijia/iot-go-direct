package bsp

import (
	"context"
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
			return
		case <-*ledManager.HeartbeatCh:
			isRefreshed = true

		case <-ticker.C:
			if !isRefreshed {
				TurnOffLed(ledManager.DeviceName)
			} else {
				TurnOnLed(ledManager.DeviceName)
				isRefreshed = false
			}
		}
	}
}
