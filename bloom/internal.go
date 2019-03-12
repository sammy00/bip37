package bloom

import (
	"github.com/sammyne/murmur3"
)

// add takes the given data into pattern record
func (f *Filter) add(data []byte) error {
	if nil == f.snapshot {
		return ErrUninitialised
	}

	for i := uint32(0); i < f.snapshot.HashFuncs; i++ {
		bitIdx := f.hash(i, data)
		//fmt.Println(bitIdx)
		// set the j(=bitIdx%8)-th bit of the k(=bitIdx/8)-th byte
		f.snapshot.Bits[bitIdx>>3] |= (1 << (bitIdx & 0x07))
	}

	return nil
}

// hash estimates the bit index mapped from the given data for the i-th
// Murmur3 employed by the filter. `idx` is used to differentiate the
// seed for each Murmur3 hashing, where the seed used by the i-th Murmur3
// would be `idx*C + Tweak`
func (f *Filter) hash(idx uint32, data []byte) uint32 {
	// seed = idx*C + f.snapshot.Tweak
	bitIdx := murmur3.SumUint32(data, idx*f.c+f.snapshot.Tweak)
	//fmt.Printf("%d-%d-%d-%x: %d\n", idx, f.c, f.snapshot.Tweak, data, bitIdx)
	return bitIdx % (uint32(len(f.snapshot.Bits)) << 3)
}

// match checks if the given data pattern is possibly recorded by the filter
func (f *Filter) match(data []byte) bool {
	if nil == f.snapshot {
		return false
	}

	// iterating each hash output and ensure the corresponding is set
	// otherwise return false
	for i := uint32(0); i < f.snapshot.HashFuncs; i++ {
		bitIdx := f.hash(i, data)
		if 0 == f.snapshot.Bits[bitIdx>>3]&(1<<(bitIdx&0x07)) {
			return false
		}
	}

	return true
}
