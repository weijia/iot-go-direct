package util

import "encoding/hex"

func DecodeHex(i string) []byte {
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
	return int(DecodeHex(string(oneChar))[0]) - 1
}
