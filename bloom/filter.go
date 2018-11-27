package bloom

import (
	"math"
	"sync"

	"github.com/sammy00/bip37/command"
)

var ln2Sqr = math.Ln2 * math.Ln2

type Filter struct {
	mtx      sync.Mutex
	snapshot *command.Load
	c        uint32
}

func (f *Filter) Add(data []byte) error {
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

func (f *Filter) Clear() {
	f.snapshot = nil
}

func (f *Filter) Loaded() bool {
	return nil == f.snapshot
}

func (f *Filter) Match(data []byte) bool {
	//return f.match(data)
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

func (f *Filter) Recover(snapshot *command.Load) *Filter {
	f.snapshot = snapshot

	return f
}

func (f *Filter) Snapshot() *command.Load {
	return f.snapshot
}

func Load(snapshot *command.Load) *Filter {
	return new(Filter).Recover(snapshot)
}

func New(N uint32, P float64, flags UpdateType,
	tweaks ...uint32) *Filter {
	P = math.Max(1e-9, math.Min(P, 1))

	// calculates S = -1/ln2Sqr*N*ln(P)/8
	S := uint32(-1 / ln2Sqr * float64(N) * math.Log(P) / 8)
	// normalize S to range (0, MaxFilterSize]
	S = MinUint32(S, MaxFilterSize)

	// calculates the nHashFuncs = S*8/N*ln2
	nHashFuncs := uint32(float64(S*8) / float64(N) * math.Ln2)
	// normalize nHashFuncs to range (0, MaxHashFuncs)
	nHashFuncs = MinUint32(nHashFuncs, MaxHashFuncs)

	c, tweak := C, Tweak
	if len(tweaks) > 2 {
		c, tweak = tweaks[0], tweaks[1]
	}

	return &Filter{
		snapshot: &command.Load{
			Bits:      make([]byte, S),
			HashFuncs: nHashFuncs,
			Tweak:     tweak,
		},
		c: c,
	}
}
