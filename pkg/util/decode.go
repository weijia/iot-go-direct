package util

import "encoding/hex"

func DecodeId(i string) []byte {
	b, err := hex.DecodeString(i)
	if err != nil {
		util.IotLogError(err)
	}
	return b
}
