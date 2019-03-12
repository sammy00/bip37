package bloom

const (
	// MaxFilterSize limits the size of the maximum length of the bit pattern
	// in bytes
	MaxFilterSize = 36000
	// MaxHashFuncs limits the maximum number of hash functions employed
	MaxHashFuncs = 50
	// C is an extra parameter tweaking the seed to initial the Murmur3.
	// See https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#bloom-filter-format
	C uint32 = 0xfba4c795
	// Tweak is the recommened random value tweaking seed for Murmur3.
	// See https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#bloom-filter-format
	Tweak uint32 = 0x00000005
)
