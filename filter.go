package bip37

import (
	"math"
)

var ln2Sqr = math.Ln2 * math.Ln2

type BloomFilter struct {
	snapshot *Snapshot
}

func (f *BloomFilter) Add(data []byte) error {
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

func (f *BloomFilter) Clear() {
	f.snapshot = nil
}

func (f *BloomFilter) Loaded() bool {
	return nil == f.snapshot
}

func (f *BloomFilter) Match(data []byte) bool {
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

func (f *BloomFilter) Recover(snapshot *Snapshot) *BloomFilter {
	f.snapshot = snapshot

	return f
}

func (f *BloomFilter) Snapshot() *Snapshot {
	return f.snapshot
}

func Load(snapshot *Snapshot) *BloomFilter {
	return new(BloomFilter).Recover(snapshot)
}

//func New(N, C, tweak uint32, P float64) *BloomFilter {
func New(N uint32, P float64, flags BloomUpdateType, 
	tweaks ...uint32) *BloomFilter {
	P = math.Max(1e-9, math.Min(P, 1))

	// calculates S = -1/ln2Sqr*N*ln(P)/8
	S := uint32(-1 / ln2Sqr * float64(N) * math.Log(P) / 8)
	// normalize S to range (0, MaxBloomFilterSize]
	S = MinUint32(S, MaxFilterSize)

	// calculates the nHashFuncs = S*8/N*ln2
	nHashFuncs := uint32(float64(S*8) / float64(N) * math.Ln2)
	// normalize nHashFuncs to range (0, MaxHashFuncs)
	nHashFuncs = MinUint32(nHashFuncs, MaxHashFuncs)

	C,tweak:=C,Tweak
	if len(tweaks)>2 {
		C, tweak = tweaks[0],tweaks[1]		
	}

	return &BloomFilter{
		snapshot: &Snapshot{
			Bits:      make([]byte, S),
			HashFuncs: nHashFuncs,
			C:         C,
			Tweak:     tweak,
		},
	}
}
