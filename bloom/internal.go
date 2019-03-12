package bloom

import (
	"github.com/sammyne/murmur3"
)

func (f *Filter) add(data []byte) error {
	if nil == f.snapshot {
		return ErrUninitialised
	}

	for i := uint32(0); i < f.snapshot.HashFuncs; i++ {
		bitIdx := f.hash(i, data)
		//fmt.Println(bitIdx)
		// set the j(=bitIdx%8)-th bit of the k()=bitIdx/8)-th byte
		f.snapshot.Bits[bitIdx>>3] |= (1 << (bitIdx & 0x07))
	}

	return nil
}

func (f *Filter) hash(idx uint32, data []byte) uint32 {
	// seed = idx*C + f.snapshot.Tweak
	bitIdx := murmur3.SumUint32(data, idx*f.c+f.snapshot.Tweak)
	//fmt.Printf("%d-%d-%d-%x: %d\n", idx, f.c, f.snapshot.Tweak, data, bitIdx)
	return bitIdx % (uint32(len(f.snapshot.Bits)) << 3)
}

func (f *Filter) match(data []byte) bool {
	if nil == f.snapshot {
		return false
	}

	for i := uint32(0); i < f.snapshot.HashFuncs; i++ {
		bitIdx := f.hash(i, data)
		if 0 == f.snapshot.Bits[bitIdx>>3]&(1<<(bitIdx&0x07)) {
			return false
		}
	}

	return true
}
