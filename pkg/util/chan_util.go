package util

import (
	"fmt"
	"time"
)

func SendMsgWithoutBlockingCommon(m interface{}, ch chan interface{}, errMsg string) {
	go func() {
		select {
		case ch <- m:
		default:
			IotLogErrorWithFormatStr(errMsg+": %v", m)
		}
	}()
}

func SendBytesMsgWithoutBlocking(m []byte, ch chan []byte, errMsg string) {
	go func() {
		select {
		case ch <- m:
		default:
			IotLogErrorWithFormatStr(errMsg+": %v", m)
		}
	}()
}

func SendRepliedNodeIdWithoutBlocking(nodeId string, ch chan<- string, timeoutSeconds int) {
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	go func() {
		select {
		case ch <- nodeId:
			eventTimer.Stop()
		case <-eventTimer.C:
			IotLogInfo(fmt.Sprintf("Timeout for notifying node reply received for node id: %s, ch: %p\n", nodeId, ch))
		}
	}()
}

var HeartbeatIndex = 0

func IsReplyTimeout(nodeIdStr string, ch <-chan string, timeoutSeconds int) bool {
	localIndex := HeartbeatIndex
	HeartbeatIndex += 1
	IotLogInfo(fmt.Sprintf("handling reply for index: %d", localIndex))
	eventTimer := time.NewTimer(time.Duration(timeoutSeconds) * time.Second)
	for {
		select {
		case <-eventTimer.C:
			IotLogInfo(fmt.Sprintf("Timeout waiting for index: %d, node id: %s", localIndex, nodeIdStr))
			// IotLogInfo(fmt.Sprintf("Timeout waiting for node id: %s", nodeIdStr))
			return true
		case responseNodeId := <-ch:
			IotLogInfo(fmt.Sprintf("Reply for node id for index: %d, node id: %s from ch: %p", localIndex, responseNodeId, ch))
			// IotLogInfo(fmt.Sprintf("Reply for node id: %s", responseNodeId))
			eventTimer.Stop()
			return false
		}
	}
}
