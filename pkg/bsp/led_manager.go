package bsp

import "context"

type LedManager struct {
	HeartbeatCh *chan<- int
	DeviceName  string
}

func (ledManager LedManager) LedMsgLoop(ctx context.Context) {

}
