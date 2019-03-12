package wire

// BloomUpdateType enumerates the updating policy for the bloom filter
type BloomUpdateType uint8

// Enumerations of different updating policies for the bloom filter
const (
	UpdateNone         BloomUpdateType = 0
	UpdateAll          BloomUpdateType = 1
	UpdateP2PubKeyOnly BloomUpdateType = 2
)
