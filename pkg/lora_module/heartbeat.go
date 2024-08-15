package lora_module

import (
	"context"
	"sync"
	"time"

	"iot_go/pkg/bsp"
	"iot_go/pkg/node"
	"iot_go/pkg/util"
)

// Will be set by top level msg loop
var TimeoutNodeIdCh *chan string

func (loraModule LoraModule) SendHeartbeatForListInLoop(ctx context.Context,
	nodeList []string) {
		
	ticker1 := time.NewTicker(time.Duration(bsp.BspConfigInstance.Heartbeat) * time.Second)
	defer ticker1.Stop()
	for {
		select{
		case nodeList = <- *loraModule.NodeListCh:
		case <- ctx.Done():
			return
		case <-ticker1.C:
			loraModule.SendHeartbeatForList(ctx, nodeList, nil)
		}
	}
}

func (loraModule LoraModule) UpdateNodeList(nodeList []string) {
	copiedNodeList := make([]string, len(nodeList)) // 预先分配足够的空间  
    copy(copiedNodeList, nodeList) // 复制元素

	select {
	case *loraModule.NodeListCh <- copiedNodeList:
	default:
		// 如果channel已满，打印消息
		util.IotLogErrorStr("UpdateNodeList: Channel full, could not send without blocking")
	}
}

func (loraModule LoraModule) SendHeartbeatForList(ctx context.Context,
	nodeList []string, wg *sync.WaitGroup) {

	module := bsp.BspConfigInstance.Module0

	if loraModule.LoraIndex == 1 {
		module = bsp.BspConfigInstance.Module1
	} else {
		if loraModule.LoraIndex == 2 {
			module = bsp.BspConfigInstance.Module2
		}
	}
	for _, nodeIdStr := range nodeList {
		util.IotDebugPrintf("Sending heartbeat on %d for %s", loraModule.LoraIndex, nodeIdStr)
		// loraModule.Mutex.Lock()
		reply := loraModule.SendNodeMsgWithRetryOrTimeout(node.GetHeartbeatMsg(nodeIdStr), util.HEARTBEAT_RETRY_CNT, node.HEARTBEAT_REPLY)
		// loraModule.Mutex.Unlock()
		// util.IotLog("SendHeartbeatForNode returned timeout & data: %v", reply)
		if reply == nil {
			util.IotLogErrorWithFormatStr("Heartbeat for %s no reply, will set node "+
				"status as offline and send node init for it on public freq", nodeIdStr)
			// Send timeout msg to main thread to update node state
			if TimeoutNodeIdCh != nil {
				
				select {
				case *TimeoutNodeIdCh <- nodeIdStr:
				default:
					// 如果channel已满，打印消息
					util.IotLogErrorStr("SendHeartbeatForList: Channel full, could not send without blocking")
				}

			}
			// Module0.Mutex.Lock()
			util.IotLog("Sending node init on 0 for %s", nodeIdStr)
			Module0.SendNodeMsgWithRetryOrTimeout(node.GetNodeInitMsg(nodeIdStr, module), util.NODE_INIT_RETRY_CNT, node.CONFIG_NODE_REPLY)
			// Module0.Mutex.Unlock()
		}
	}
	if wg != nil {
		wg.Done()
	}
}
