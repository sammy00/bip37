package bloom

import (
	btcwire "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/sammyne/bip37/wire"
)

// AddOutPoint takes a COutPoint into record, which is actually the
// concurrent-safe version of addOutPoint
func (f *Filter) AddOutPoint(out *btcwire.OutPoint) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.addOutPoint(out.Hash[:], out.Index)
}

// MatchOutPoint checks if the given COutPoint is possibly recorded
// in the bit pattern of the filter
func (f *Filter) MatchOutPoint(out *btcwire.OutPoint) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.match(marshalOutPoint(out))
}

// MatchTx checks if the tx matches the bit pattern of filter
func (f *Filter) MatchTx(tx *btcutil.Tx) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	bak := f.snapshot.Flags
	f.snapshot.Flags = wire.UpdateNone
	ok := f.matchTxAndUpdate(tx)
	f.snapshot.Flags = bak

	return ok
}

// MatchTxAndUpdate checks if the tx matches the bit pattern of filter and
// update the bit pattern accordingly
func (f *Filter) MatchTxAndUpdate(tx *btcutil.Tx) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.matchTxAndUpdate(tx)
}
