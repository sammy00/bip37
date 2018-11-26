package bloom

type Snapshot struct {
	Bits      []byte
	HashFuncs uint32
	C         uint32
	Tweak     uint32
	//Flags     UpdateType
}
