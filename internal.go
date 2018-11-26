package bloom

import (
	"github.com/sammy00/murmur3"
)

/*
func (f *Filter) add(data []byte) error {
	if nil == f.snapshot {
		return ErrUninitialised
	}

	for i := uint32(0); i < f.snapshot.HashFuncs; i++ {
		bitIdx := f.hash(i, data)
		// set the j(=bitIdx%8)-th bit of the k()=bitIdx/8)-th byte
		f.snapshot.Bits[bitIdx>>3] |= (1 << (bitIdx & 0x0f))
	}

	return nil
}
*/

func (f *Filter) hash(idx uint32, data []byte) uint32 {
	// seed = idx*C + f.snapshot.Tweak
	bitIdx := murmur3.SumUint32(data, idx*f.snapshot.C+f.snapshot.Tweak)
	return bitIdx % (uint32(len(f.snapshot.Bits)) << 3)
}

/*
func (f *Filter) match(data []byte) bool {
	if nil == f.snapshot {
		return false
	}

	for i := uint32(0); i < f.snapshot.HashFuncs; i++ {
		bitIdx := f.hash(i, data)
		if 0 == f.snapshot.Bits[bitIdx>>3]&(1<<(bitIdx&0x0f)) {
			return false
		}
	}

	return true
}
*/
