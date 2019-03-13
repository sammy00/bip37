package bloom

import (
	"math"
	"sync"

	"github.com/sammyne/bip37/wire"
)

var ln2Sqr = math.Ln2 * math.Ln2

// Filter implements a concurrent safe bloom filter
type Filter struct {
	mtx      sync.Mutex
	snapshot *wire.FilterLoad
	c        uint32
}

// Add is the concurrently safe version of its unexported variant `add`
func (f *Filter) Add(data []byte) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.add(data)
}

// Clear resets the filter by empty its bit pattern, which is safe for
// concurrent use
func (f *Filter) Clear() {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.snapshot = nil
}

// Loaded checks if the filter has been initialized properly, which is safe for
// concurrent use
func (f *Filter) Loaded() bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return nil != f.snapshot
}

// Match checks if the data may be recorded by the filter, which is safe for
// concurrent use
func (f *Filter) Match(data []byte) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.match(data)
}

// Recover overrides the bit pattern in a concurrently safe manner
func (f *Filter) Recover(snapshot *wire.FilterLoad) *Filter {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.snapshot = snapshot

	return f
}

// Snapshot return the bit pattern maintained by filter up till now
func (f *Filter) Snapshot() *wire.FilterLoad {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.snapshot
}

// Load is the procedural version of Filter.Recover
func Load(snapshot *wire.FilterLoad) *Filter {
	return new(Filter).Recover(snapshot)
}

// New serves as the constructor of a bloom filter according to
// specification in https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#bloom-filter-format
func New(N uint32, P float64, flags wire.BloomUpdateType,
	tweaks ...uint32) *Filter {
	// false positive rate
	P = math.Max(1e-9, math.Min(P, 1))

	// calculates size of the filter S = -1/ln2Sqr*N*ln(P)/8
	S := uint32(-1 / ln2Sqr * float64(N) * math.Log(P) / 8)
	// normalize S to range (0, MaxFilterSize]
	S = MinUint32(S, MaxFilterSize)

	// calculates the nHashFuncs = S*8/N*ln2
	nHashFuncs := uint32(float64(S*8) / float64(N) * math.Ln2)
	// normalize nHashFuncs to range (0, MaxHashFuncs)
	nHashFuncs = MinUint32(nHashFuncs, MaxHashFuncs)

	c, tweak := C, Tweak
	if len(tweaks) >= 1 {
		tweak = tweaks[0]
	}
	if len(tweaks) >= 2 {
		c = tweaks[1]
	}

	return &Filter{
		snapshot: &wire.FilterLoad{
			Bits:      make([]byte, S),
			HashFuncs: nHashFuncs,
			Tweak:     tweak,
			Flags:     flags,
		},
		c: c,
	}
}
