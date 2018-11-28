package wire

type BloomUpdateType uint8

const (
	UpdateNone         BloomUpdateType = 0
	UpdateAll          BloomUpdateType = 1
	UpdateP2PubKeyOnly BloomUpdateType = 2
)
