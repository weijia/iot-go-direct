package msg

type GroupUpdateGlassColorRequest struct {
	Method string                        `json:"method"`
	Params []GroupUpdateGlassColorParams `json:"params"`
}
type GroupUpdateGlassColorParams struct {
	NodeList1 []string `json:"node_list1"`
	NodeList2 []string `json:"node_list2"`
	Color     string   `json:"color"`
}
