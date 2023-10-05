package node

import (
	"fmt"
	"iot_go/pkg/bsp"
	"iot_go/pkg/shared"
	"iot_go/pkg/util"
	"time"
)

type NodeMsgBase struct {
	MsgLen    int     `struc:"uint8"`
	BoardType int     `struc:"uint8"`
	CmdType   int     `struc:"uint8"`
	GatewayId []byte  `struc:"[6]uint8"`
	NodeId    [4]byte `struc:"[4]uint8"`
}

const (
	HEARTBEAT_REQ            = 1
	HEARTBEAT_REPLY          = 2
	CONFIG_NODE_REQ          = 3
	CONFIG_NODE_REPLY        = 4
	UPDATE_GLASS_COLOR_REQ   = 5
	UPDATE_GLASS_COLOR_REPLY = 6
	GET_GLASS_STATE_REQ      = 9
	GET_GLASS_STATE_REPLY    = 10

	REPLY_CMD_START_INDEX                   = 2
	REPLY_GATEWAY_ID_START_INDEX            = 3
	GATEWAY_ID_LEN                          = 6
	REPLY_NODE_ID_START_INDEX               = REPLY_GATEWAY_ID_START_INDEX + GATEWAY_ID_LEN // 9
	NODE_ID_LEN                             = 4
	REPLY_PAYLOAD_START_INDEX               = REPLY_NODE_ID_START_INDEX + NODE_ID_LEN // 13
	UPDATE_GLASS_COLOR_RESULT_INDEX         = REPLY_PAYLOAD_START_INDEX // 13
	HEARTBEAT_COLOR_POS_START_INDEX         = REPLY_NODE_ID_START_INDEX + NODE_ID_LEN // 13
	NODE_INIT_REPLY_PARAM_LEN               = 3
	NODE_INIT_REPLY_PARAM_TYPE_FREQ         = 1
	NODE_INIT_REPLY_PARAM_TYPE_FACTOR       = 2
	NODE_INIT_REPLY_PARAM_TYPE_BAND         = 3
	NODE_INIT_REPLY_PARAM_TYPE_HW_VER       = 4
	NODE_INIT_REPLY_PARAM_TYPE_SW_VER       = 5
	NODE_INIT_REPLY_PARAM_TYPE_RUNNING_AREA = 6
	// TODO: Need to adjust this timer
	NODE_INIT_REPLY_TIMEOUT_SECONDS = 6

	UPDATE_COLOR_RESULT_OK   = 1
	UPDATE_COLOR_RESULT_FAIL = 0
)

func GetNodeIdStr(msg []byte) string {
	return fmt.Sprintf("%x", msg[REPLY_NODE_ID_START_INDEX:REPLY_NODE_ID_START_INDEX+NODE_ID_LEN])
}

// type NodeInitReplyParam struct {
// 	ParamType int `struc:"uint8"`
// 	ParamValue int `struc:"uint16"`
// }

// type NodeInitReply struct {
// 	NodeMsgBase
// 	ParamNum int `struc:"uint8,sizeof=Str"`
//     ParamList   []NodeInitReplyParam `struc:"NodeInitReplyParam,sizeof=Str"`
// }

func GetNodeIdFromMsg(reply []byte) string {
	return fmt.Sprintf("%x", reply[REPLY_NODE_ID_START_INDEX:REPLY_NODE_ID_START_INDEX+NODE_ID_LEN])
}

// We can only handle 1 config request at a time
var OngoingConfigReqParam shared.ConfigParams
var IsProcessingConfigReq bool = false

// var PendingInitReqNodeList []string

var ConfigReqCh = make(chan shared.ConfigParams)

// var NodeConfigTimer = time.NewTimer(NODE_INIT_REPLY_TIMEOUT_SECONDS * time.Second)
// var CancelFuncForNodeInitReplyTimeout context.CancelFunc

func UpdateNodeStateForInitReply(nodeInitReply []byte) {
	util.IotLogInfo(fmt.Sprintf("Got init reply: %v\n", nodeInitReply))
	nodeId := GetNodeIdFromMsg(nodeInitReply)

	// Update node state in bsp
	nodeState := bsp.GetOrCreateNodeState(nodeId)
	if nodeState == nil {
		nodeState = &bsp.NodeState{}
		bsp.BspConfigInstance.NodeStates = append(bsp.BspConfigInstance.NodeStates, *nodeState)
	}

	paramNum := int(nodeInitReply[REPLY_PAYLOAD_START_INDEX])

	for i := 0; i < paramNum; i++ {
		paramType := int(nodeInitReply[REPLY_PAYLOAD_START_INDEX+1+i*NODE_INIT_REPLY_PARAM_LEN])
		paramValue := int(nodeInitReply[REPLY_PAYLOAD_START_INDEX+1+i*NODE_INIT_REPLY_PARAM_LEN+1])<<8 |
			int(nodeInitReply[REPLY_PAYLOAD_START_INDEX+1+i*NODE_INIT_REPLY_PARAM_LEN+2])
		switch paramType {
		case NODE_INIT_REPLY_PARAM_TYPE_FREQ:
			nodeState.ModuleParam.Freq = paramValue
		case NODE_INIT_REPLY_PARAM_TYPE_FACTOR:
			nodeState.ModuleParam.Factor = paramValue
		case NODE_INIT_REPLY_PARAM_TYPE_BAND:
			nodeState.ModuleParam.Band = paramValue
		case NODE_INIT_REPLY_PARAM_TYPE_HW_VER:
			nodeState.HwVer = paramValue
		case NODE_INIT_REPLY_PARAM_TYPE_SW_VER:
			nodeState.SwVer = paramValue
		case NODE_INIT_REPLY_PARAM_TYPE_RUNNING_AREA:
			nodeState.RunningArea = paramValue
		}
	}
	nodeState.LastMsgTimestamp = time.Now().Unix()

	// for i, NodeInitPendingNodeId := range PendingInitReqNodeList {
	// 	if NodeInitPendingNodeId == nodeId {
	// 		PendingInitReqNodeList = append(PendingInitReqNodeList[:i], PendingInitReqNodeList[i+1:]...)
	// 		break
	// 	}
	// }
	// if len(PendingInitReqNodeList) == 0 {
	// 	CancelFuncForNodeInitReplyTimeout()
	// 	util.IotLogInfo(fmt.Sprintf("config request completed, timer canceled, param point: %p\n", &OngoingConfigReqParam))
	// 	ConfigReqCh <- OngoingConfigReqParam
	// }
}
