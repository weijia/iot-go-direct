package msg

import (
	"encoding/json"
	"fmt"
	"testing"

	"iot_go/pkg/msg"
)

func TestJsonUnmarshal(t *testing.T) {

	init := &msg.Init{
		MsgType:       "init",
		GatewayNodeID: "testing",
	}
	init_json_str, _ := json.Marshal(init)
	fmt.Println(string(init_json_str))
	b := []byte(init_json_str)
	var dat msg.BaseMsg
	if err := json.Unmarshal(b, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	// fmt.Println("hello test")
}
