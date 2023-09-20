package util

import (
	"fmt"
	"time"
)

func SendRepliedNodeIdWithoutBlocking(nodeId string, ch chan<- string, timeoutSeconds int) {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	go func() {
		select {
		case ch <- nodeId:
			eventTimer.Stop()
		case <-eventTimer.C:
			IotLogInfo(fmt.Sprintf("Timeout sending for node id: %s\n", nodeId))
		}
	}()
}

func IsReplyTimeout(nodeIdStr string, ch <-chan string, timeoutSeconds int) bool {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	for {
		select {
		case <-eventTimer.C:
			IotLogInfo(fmt.Sprintf("Timeout waiting for node id: %s\n", nodeIdStr))
			return true
		case responseNodeId := <-ch:
			IotLogInfo(fmt.Sprintf("Reply for node id: %s\n", responseNodeId))
			eventTimer.Stop()
			return false
		}
	}
}
