package bloom_test

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	"github.com/btcsuite/btcutil"
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
		filter := bloom.New(3, 0.01, wire.UpdateAll, Tweak, C)

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

func TestFilter_Add_pubKey(t *testing.T) {
	priv := "5Kg1gnAjaLfKiwhhPpGS3QfRg2m6awQvaj98JCZBZQ5SuS2F15C"

	wif, err := btcutil.DecodeWIF(priv)
	if err != nil {
		t.Errorf("TestFilterInsertKey DecodeWIF failed: %v", err)
		return
	}

	f := bloom.New(2, 0.001, wire.UpdateAll, 0)
	f.Add(wif.SerializePubKey())
	f.Add(btcutil.Hash160(wif.SerializePubKey()))

	expect := &wire.FilterLoad{
		Bits:      bloom.Hexlify("8fc16b"),
		HashFuncs: 8,
		Tweak:     0,
		Flags:     1,
	}

	if got := f.Snapshot(); !reflect.DeepEqual(got, expect) {
		t.Fatalf("invalid snapshot: got %v, expect %v", got, expect)
	}
}

func TestFilter_Clear(t *testing.T) {
	filter := bloom.New(8, 0.123, wire.UpdateAll)
	filter.Clear()

	if got := filter.Snapshot(); got != nil {
		t.Fatalf("cleared filter should produce nil snapshot")
	}
}

func TestFilter_Load(t *testing.T) {
	expect := &wire.FilterLoad{Bits: []byte("hello world")}

	if got := bloom.Load(expect).Snapshot(); got != expect {
		t.Fatal("failed to load filter from snapshot")
	}
}

func TestFilter_Loaded(t *testing.T) {
	testCases := []struct {
		clear  bool
		expect bool
	}{
		{false, true},
		{true, false},
	}

	for i, c := range testCases {
		filter := bloom.New(8, 0.123, wire.UpdateAll)
		if c.clear {
			filter.Clear()
		}

		if got := filter.Loaded(); got != c.expect {
			t.Fatalf("#%d unexpected loading status: got %v, expect %v", i, got,
				c.expect)
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
		filter := bloom.New(3, 0.01, wire.UpdateAll, Tweak, C)
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

func TestFilter_Recover(t *testing.T) {
	snapshot := &wire.FilterLoad{Bits: []byte("hello world")}

	filter := new(bloom.Filter).Recover(snapshot)

	if got := filter.Snapshot(); got != snapshot {
		t.Fatal("snapshot isn't recovered correctly")
	}
}

func TestNew(t *testing.T) {
	type expect struct {
		bitsLen    int
		nHashFuncs uint32
	}

	testCases := []struct {
		description string
		N           uint32
		P           float64 // false positive rate
		expect      expect
	}{
		{
			"fprates > 1 should be clipped at 1",
			1, 20.9999999769,
			expect{0, 0},
		},
		{
			"fprates less than 1e-9 should be clipped at min",
			1, 0,
			expect{5, 27},
		},
		{
			"negative fprates should be clipped at min",
			1, -1,
			expect{5, 27},
		},
		{
			"fprates > 1 should be clipped at 1 #2",
			8, math.E * math.Ln2 * math.Ln2 * 2,
			expect{0, 0},
		},
	}

	for i, c := range testCases[3:] {
		filter := bloom.New(c.N, c.P, wire.UpdateAll)

		snapshot := filter.Snapshot()
		bitsLen, nHashFuncs := len(snapshot.Bits), snapshot.HashFuncs

		if bitsLen != c.expect.bitsLen {
			t.Fatalf("#%d %s: got %d, expect %d", i, c.description,
				bitsLen, c.expect.bitsLen)
		}
		if nHashFuncs != c.expect.nHashFuncs {
			t.Fatalf("#%d %s: got %d, expect %d", i, c.description,
				nHashFuncs, c.expect.nHashFuncs)
		}
	}
}

func TestNew_withTweak(t *testing.T) {
	testCases := []struct {
		filter *bloom.Filter
		data   [][]byte
		expect *wire.FilterLoad
	}{
		{
			bloom.New(3, 0.01, wire.UpdateAll, 2147483649),
			[][]byte{
				bloom.Hexlify("99108ad8ed9bb6274d3980bab5a85c048f0950c8"),
				//bloom.Hexlify("19108ad8ed9bb6274d3980bab5a85c048f0950c8"),
				bloom.Hexlify("b5a2c786d9ef4658287ced5914b37a1b4aa32eee"),
				bloom.Hexlify("b9300670b4c5366e95b2699e8b18bc75e5f729c5"),
			},
			&wire.FilterLoad{
				Bits:      bloom.Hexlify("ce4299"),
				HashFuncs: 5,
				Tweak:     2147483649,
				Flags:     wire.UpdateAll,
			},
		},
	}

	for i, c := range testCases {
		for _, v := range c.data {
			c.filter.Add(v)
		}

		if got := c.filter.Snapshot(); !reflect.DeepEqual(got, c.expect) {
			t.Fatalf("#%d invalid snapshot: got %v, expect %v", i, got, c.expect)
		}
	}
}
