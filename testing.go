package bip37

import (
	"encoding/hex"
)

func Hexlify(str string) []byte {
	out, _ := hex.DecodeString(str)
	return out
}
