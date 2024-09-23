package bsp

import (
	"encoding/hex"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
)

type NodeState struct {
	NodeId              string        `json:"node_id"`
	LastMsgTimestamp    int64         `json:"node_msg_timestamp"`
	NodeReportedColor   [8]int        `json:"node_reported_color"`
	NodeRequestingColor [8]int        `json:"node_requesting_color"`
	HwVer               int           `json:"hardware_version"`
	SwVer               int           `json:"software_version"`
	RunningArea         int           `json:"running_area"`
	ModuleParam         shared.Module `json:"module_param"`
	RSSI                float64       `json:"rssi"`
	SNR                 float64       `json:"snr"`
	IsOffline           bool          `json:"is_offline"`
	CompletionStatus    int           `json:"completion_status"`
}

func GetOrCreateNodeState(nodeId string) *NodeState {
	for index := range BspConfigInstance.NodeStates {
		if BspConfigInstance.NodeStates[index].NodeId == nodeId {
			return &BspConfigInstance.NodeStates[index]
		}
	}
	newNode := NodeState{
		NodeId: nodeId,
	}
	BspConfigInstance.NodeStates = append(BspConfigInstance.NodeStates, newNode)
	return &BspConfigInstance.NodeStates[len(BspConfigInstance.NodeStates)-1]
}

func SetRequestingColor(nodeId string, requestingColor string) {
	nodeState := GetOrCreateNodeState(nodeId)
	SetRequestingColorByNodeStateRef(nodeState, requestingColor)
}

func SetRequestingColorByNodeStateRef(nodeState *NodeState, requestingColor string) {
	// 将十六进制字符串解码为字节切片
	bytes, err := hex.DecodeString(requestingColor)
	if err != nil {
		util.IotLogError(err)
		return
	}

	// 处理每个字节
	for _, b := range bytes {
		// 这里只是简单地将字节转换为大写（实际上字节没有大小写之分，这里仅作示例）
		// 但你可以在这里执行任何你需要的字节处理操作
		// fmt.Printf("原始字节[%d]: %02x\n", i, b)

		// 假设的“处理”：比如我们可以将字节与另一个字节进行异或操作
		// 这里只是简单地将字节与自身异或，作为示例
		// processedByte := b ^ b // 实际上这不会改变字节的值

		// 打印处理后的字节（在这个例子中，它仍然是相同的）
		// fmt.Printf("处理后的字节[%d]: %02x\n", i, processedByte)
		nodeState.NodeRequestingColor[(b&0xf0)>>4] = int(b & 0xf)
	}
}

func SetAllRequestingColor(requestingColor string) {
	for index := range BspConfigInstance.NodeStates {
		SetRequestingColorByNodeStateRef(&BspConfigInstance.NodeStates[index], requestingColor)
	}
}
