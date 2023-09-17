package msg

type GroupUpdateGlassColorReply struct {
	GatewayNodeIdReply
	InvalidNodes []string `json:"invalid_nodes"`
}