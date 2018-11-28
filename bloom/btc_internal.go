package bloom

import (
	"encoding/binary"

	"github.com/sammy00/bip37/wire"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcutil"
)

func (f *Filter) addOutPoint(txHash []byte, index uint32) error {
	var i [4]byte
	binary.LittleEndian.PutUint32(i[:], index)

	return f.add(append(txHash, i[:]...))
}

//func (f *Filter) maybeAddOutPoint(scriptClass, )

func (f *Filter) matchTxAndUpdate(tx *btcutil.Tx) bool {
	// check tx hash
	txHash := tx.Hash()[:]
	ok := f.match(txHash)

	// check elements in public key script of tx output
	for idx, out := range tx.MsgTx().TxOut {
		data, err := txscript.PushedData(out.PkScript)
		if nil != err {
			continue // skip the unexpected pushed data
		}

		for _, elem := range data {
			if !f.match(elem) {
				continue // skip the negative
			}

			ok = true
			// add the OutPoint as specified
			switch f.snapshot.Flags {
			case wire.UpdateAll:
				f.addOutPoint(txHash, uint32(idx))
			case wire.UpdateP2PubKeyOnly:
				if C := txscript.GetScriptClass(
					out.PkScript); txscript.PubKeyTy == C || txscript.MultiSigTy == C {
					f.addOutPoint(txHash, uint32(idx))
				}
			}
		}
	}

	// return if match found
	if ok {
		return true
	}

	// check OutPoint corresponding to tx input
	for _, in := range tx.MsgTx().TxIn {
		if f.match(serializeOutPoint(&in.PreviousOutPoint)) {
			return true
		}

		data, err := txscript.PushedData(in.SignatureScript)
		if nil != err {
			continue
		}

		for _, elem := range data {
			if f.match(elem) {
				return true
			}
		}
	}

	return false
}
