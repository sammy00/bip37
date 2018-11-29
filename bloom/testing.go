package bloom

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	btcwire "github.com/btcsuite/btcd/wire"
)

func Hexlify(str string) []byte {
	out, _ := hex.DecodeString(str)
	return out
}

func NewOutPoint(hash []byte, index uint32) *btcwire.OutPoint {
	var chash chainhash.Hash
	copy(chash[:], hash)

	return btcwire.NewOutPoint(&chash, index)
}

func Unhexlify(str string) []byte {
	out, _ := hex.DecodeString(str)
	return out
}

/*
func FakeTx(t *testing.T) *btcutil.Tx {
	txStr := "01000000010b26e9b7735eb6aabdf358bab62f9816a21ba9ebdb719d5299e" +
		"88607d722c190000000008b4830450220070aca44506c5cef3a16ed519d7" +
		"c3c39f8aab192c4e1c90d065f37b8a4af6141022100a8e160b856c2d43d2" +
		"7d8fba71e5aef6405b8643ac4cb7cb3c462aced7f14711a0141046d11fee" +
		"51b0e60666d5049a9101a72741df480b96ee26488a4d3466b95c9a40ac5e" +
		"eef87e10a5cd336c19a84565f80fa6c547957b7700ff4dfbdefe76036c33" +
		"9ffffffff021bff3d11000000001976a91404943fdd508053c75000106d3" +
		"bc6e2754dbcff1988ac2f15de00000000001976a914a266436d296554760" +
		"8b9e15d9032a7b9d64fa43188ac00000000"
	strBytes, err := hex.DecodeString(txStr)
	if err != nil {
		t.Fatalf("TestFilterBloomMatch DecodeString failure: %v", err)
	}

	tx, err := btcutil.NewTxFromBytes(strBytes)
	if err != nil {
		t.Errorf("TestFilterBloomMatch NewTxFromBytes failure: %v", err)
	}

	return tx
}*/

func FakeTx(t *testing.T) *btcutil.Tx {
	const txJSON = `{
		"Version": 1,
		"TxIn": [
			{
				"PreviousOutPoint": {
					"Hash": [ 11, 38, 233, 183, 115, 94, 182, 170, 
						189, 243, 88, 186, 182, 47, 152, 22, 
						162, 27, 169, 235, 219, 113, 157, 82,
						153, 232, 134, 7, 215, 34, 193, 144
					],
					"Index": 0
				},
				"SignatureScript": "SDBFAiAHCspEUGxc7zoW7VGdfDw5+KqxksThyQ0GXze4pK9hQQIhAKjhYLhWwtQ9J9j7px5a72QFuGQ6xMt8s8RirO1/FHEaAUEEbRH+5RsOYGZtUEmpEBpydB30gLlu4mSIpNNGa5XJpArF7u+H4Qpc0zbBmoRWX4D6bFR5V7dwD/Tfve/nYDbDOQ==",
				"Witness": null,
				"Sequence": 4294967295
			}
		],
		"TxOut": [
			{
				"Value": 289275675,
				"PkScript": "dqkUBJQ/3VCAU8dQABBtO8bidU28/xmIrA=="
			},
			{
				"Value": 14554415,
				"PkScript": "dqkUomZDbSllVHYIueFdkDKnudZPpDGIrA=="
			}
		],
		"LockTime": 0
	}`

	var msgTx btcwire.MsgTx
	if err := json.Unmarshal([]byte(txJSON), &msgTx); nil != err {
		t.Fatal(err)
	}

	return btcutil.NewTx(&msgTx)
}
