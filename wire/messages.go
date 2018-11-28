package wire

type FilterAdd struct {
	Data []byte
}

type FilterLoad struct {
	Bits      []byte
	HashFuncs uint32
	Tweak     uint32
	Flags     BloomUpdateType
}

type FilterClear struct{}
