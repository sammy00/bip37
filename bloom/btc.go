package bloom

import (
	btcwire "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/sammyne/bip37/wire"
)

func (f *Filter) AddOutPoint(out *btcwire.OutPoint) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.addOutPoint(out.Hash[:], out.Index)
}

func (f *Filter) MatchOutPoint(out *btcwire.OutPoint) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.match(marshalOutPoint(out))
}

func (f *Filter) MatchTx(tx *btcutil.Tx) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	bak := f.snapshot.Flags
	f.snapshot.Flags = wire.UpdateNone
	ok := f.matchTxAndUpdate(tx)
	f.snapshot.Flags = bak

	return ok
}

func (f *Filter) MatchTxAndUpdate(tx *btcutil.Tx) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.matchTxAndUpdate(tx)
}
