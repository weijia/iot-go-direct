package msg

type UpdateGlassColorReply struct {
	Method string                   `json:"method"`
	Params []UpdateGlassColorParams `json:"params"`
}
type UpdateGlassColorParams struct {
	NodeID string `json:"node_id"`
	Color  string `json:"color"`
}


func (request UpdateGlassColorRequest) handle() interface{} {
    reply.GatewayNodeID = bsp.BspConfigInstance.InitConfig.GatewayNodeID
    reply.NodeList1 = bsp.BspConfigInstance.InitConfig.NodeList1
    reply.NodeList2 = bsp.BspConfigInstance.InitConfig.NodeList2
    reply.TouchNodeList1 = bsp.BspConfigInstance.InitConfig.TouchNodeList1
    reply.TouchNodeList2 = bsp.BspConfigInstance.InitConfig.TouchNodeList2
    return reply