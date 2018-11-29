package bloom

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	btcwire "github.com/btcsuite/btcd/wire"
)

func Hexlify(str string) []byte {
	out, _ := hex.DecodeString(str)
	return out
}

func NewOutPoint(hash []byte, index uint32) *btcwire.OutPoint {
	var chash chainhash.Hash
	copy(chash[:], hash)

	return btcwire.NewOutPoint(&chash, index)
}

func Unhexlify(str string) []byte {
	out, _ := hex.DecodeString(str)
	return out
}
