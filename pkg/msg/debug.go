package msg

import (
	"encoding/json"
	"fmt"
)

func DumpMsg(msg interface{}) {
	res, _ := json.MarshalIndent(msg, "", "  ")
	fmt.Println(string(res))
}
