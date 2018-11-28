package bloom

import (
	"math"
	"sync"

	"github.com/sammy00/bip37/wire"
)

var ln2Sqr = math.Ln2 * math.Ln2

type Filter struct {
	mtx      sync.Mutex
	snapshot *wire.FilterLoad
	c        uint32
}

func (f *Filter) Add(data []byte) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.add(data)
}

func (f *Filter) Clear() {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.snapshot = nil
}

func (f *Filter) Loaded() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return nil != f.snapshot
}

func (f *Filter) Match(data []byte) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.match(data)
}

func (f *Filter) Recover(snapshot *wire.FilterLoad) *Filter {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.snapshot = snapshot

	return f
}

func (f *Filter) Snapshot() *wire.FilterLoad {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.snapshot
}

func Load(snapshot *wire.FilterLoad) *Filter {
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
		snapshot: &wire.FilterLoad{
			Bits:      make([]byte, S),
			HashFuncs: nHashFuncs,
			Tweak:     tweak,
		},
		c: c,
	}
}
