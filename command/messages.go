package command

type Add struct {
	Data []byte
}

type Load struct {
	Bits      []byte
	HashFuncs uint32
	Tweak     uint32
	Flags     uint8
}

type Clear struct{}
