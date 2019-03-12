package bloom

import (
	"encoding/binary"

	"github.com/btcsuite/btcd/wire"
)

// marshalOutPoint marshals a tx output interpreted as point as `hash||index`,
// where the index is encoded in little-endian
func marshalOutPoint(out *wire.OutPoint) []byte {
	var i [4]byte
	binary.LittleEndian.PutUint32(i[:], out.Index)

	return append(out.Hash[:], i[:]...)
}
