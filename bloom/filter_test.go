package bloom_test

import (
	"bytes"
	"testing"

	"github.com/sammy00/bip37/bloom"
	"github.com/sammy00/bip37/wire"
)

const (
	Tweak = 0x00000005
	C     = 0xfba4c795
)

func TestFilter_Add(t *testing.T) {
	testCases := []struct {
		data   []byte
		expect []byte // the expected bits array in snapshot
	}{
		{
			bloom.Hexlify("99108ad8ed9bb6274d3980bab5a85c048f0950c8"),
			bloom.Hexlify("0021c1"),
		},
		{
			bloom.Hexlify("19108ad8ed9bb6274d3980bab5a85c048f0950c8"),
			bloom.Hexlify("00202e"),
		},
		{
			bloom.Hexlify("b5a2c786d9ef4658287ced5914b37a1b4aa32eee"),
			bloom.Hexlify("064c00"),
		},
		{
			bloom.Hexlify("b9300670b4c5366e95b2699e8b18bc75e5f729c5"),
			bloom.Hexlify("148402"),
		},
	}

	for i, c := range testCases {
		filter := bloom.New(3, 0.01, wire.UpdateAll, C, Tweak)

		if err := filter.Add(c.data); nil != err {
			t.Fatalf("#%d unexpected error: %v", i, err)
		}

		//btc := btcbloom.NewFilter(3, Tweak, 0.01, wire.BloomUpdateAll)
		//btc.Add(c.data)
		//c.expect = btc.MsgFilterLoad().Filter
		//t.Logf(`bloom.Hexlify("%x")`, c.expect)

		if got := filter.Snapshot().Bits; !bytes.Equal(got, c.expect) {
			t.Fatalf("#%d invalid bits: got %v, expect %v", i, got, c.expect)
		}
	}
}

func TestFilter_Match(t *testing.T) {
	testCases := []struct {
		data  []byte
		added bool
	}{
		{bloom.Hexlify("99108ad8ed9bb6274d3980bab5a85c048f0950c8"), true},
		{bloom.Hexlify("19108ad8ed9bb6274d3980bab5a85c048f0950c8"), false},
		{bloom.Hexlify("b5a2c786d9ef4658287ced5914b37a1b4aa32eee"), false},
		{bloom.Hexlify("b9300670b4c5366e95b2699e8b18bc75e5f729c5"), true},
	}

	for i, c := range testCases {
		filter := bloom.New(3, 0.01, wire.UpdateAll, C, Tweak)
		if c.added {
			filter.Add(c.data)
		}

		if c.added && !filter.Match(c.data) {
			t.Fatalf("#%d failed to match added data: %x", i, c.data)
		} else if !c.added && filter.Match(c.data) {
			t.Fatalf("#%d unexpected false positive: %x", i, c.data)
		}
	}
}
