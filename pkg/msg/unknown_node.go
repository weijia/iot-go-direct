package msg

type UnknownNode struct {
	MsgType string `json:"msg_type"`
	UnknownNodeList      []string `json:"unknown_node_list"`
}