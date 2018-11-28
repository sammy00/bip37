package bloom

import (
	"encoding/binary"

	"github.com/btcsuite/btcd/wire"
)

func serializeOutPoint(out *wire.OutPoint) []byte {
	var i [4]byte
	binary.LittleEndian.PutUint32(i[:], out.Index)

	return append(out.Hash[:], i[:]...)
}
