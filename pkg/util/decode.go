package util

import "encoding/hex"

func DecodeHex(i string) []byte {
	if len(i)%2 == 1 {
		i = "0" + i
	}
	b, err := hex.DecodeString(i)
	if err != nil {
		IotLogError(err)
	}
	return b
}

func DecodeId(i string) []byte {
	return DecodeHex(i)
}

func GetGlassAreaFromStr(oneChar byte) int {
	// glass area start from 1 so change to internal index need to -1
	return int(DecodeHex(string(oneChar))[0]) - 1
}
