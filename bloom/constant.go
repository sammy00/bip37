package bloom

type UpdateType uint8

const (
	None         UpdateType = 0
	All          UpdateType = 1
	P2PubKeyOnly            = 2
)

const (
	MaxFilterSize        = 36000
	MaxHashFuncs         = 50
	C             uint32 = 0xfba4c795
	Tweak         uint32 = 0x00000005
)
