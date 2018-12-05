package merkle_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sammy00/bip37/bloom"
	"github.com/sammy00/bip37/merkle"

	"github.com/sammy00/bip37"
	"github.com/sammy00/bip37/wire"

	btcwire "github.com/btcsuite/btcd/wire"
)

func TestNew(t *testing.T) {
	const expectJSON = `
	{
		"Header": {
			"Version": 1,
			"PrevBlock": [
				130,
				187,
				134,
				156,
				243,
				167,
				147,
				67,
				42,
				102,
				232,
				38,
				224,
				90,
				111,
				195,
				116,
				105,
				248,
				239,
				183,
				66,
				29,
				200,
				128,
				103,
				1,
				0,
				0,
				0,
				0,
				0
			],
			"MerkleRoot": [
				127,
				22,
				197,
				150,
				46,
				139,
				217,
				99,
				101,
				156,
				121,
				60,
				227,
				112,
				217,
				95,
				9,
				59,
				199,
				227,
				103,
				17,
				123,
				60,
				48,
				193,
				248,
				253,
				208,
				217,
				114,
				135
			],
			"Timestamp": "2010-12-29T21:32:38+08:00",
			"Bits": 453281356,
			"Nonce": 696601429
		},
		"Transactions": 7,
		"Hashes": [
			[
				11,
				54,
				116,
				198,
				229,
				15,
				54,
				243,
				111,
				122,
				159,
				72,
				94,
				118,
				199,
				134,
				139,
				244,
				217,
				245,
				152,
				78,
				170,
				11,
				89,
				150,
				101,
				120,
				118,
				170,
				124,
				20
			],
			[
				253,
				172,
				249,
				179,
				235,
				7,
				116,
				18,
				231,
				169,
				104,
				210,
				228,
				241,
				27,
				154,
				157,
				238,
				49,
				45,
				102,
				97,
				135,
				237,
				119,
				238,
				125,
				38,
				175,
				22,
				203,
				11
			],
			[
				138,
				146,
				163,
				234,
				16,
				184,
				199,
				40,
				160,
				183,
				225,
				11,
				57,
				177,
				177,
				226,
				129,
				214,
				72,
				157,
				90,
				55,
				22,
				242,
				40,
				38,
				142,
				129,
				87,
				134,
				115,
				72
			],
			[
				65,
				192,
				91,
				223,
				113,
				100,
				50,
				103,
				222,
				210,
				207,
				3,
				122,
				242,
				16,
				90,
				3,
				102,
				33,
				252,
				244,
				104,
				88,
				188,
				29,
				72,
				240,
				82,
				160,
				31,
				152,
				2
			],
			[
				1,
				159,
				91,
				1,
				212,
				25,
				94,
				203,
				201,
				57,
				143,
				191,
				60,
				59,
				31,
				169,
				187,
				49,
				131,
				48,
				29,
				122,
				31,
				179,
				189,
				23,
				79,
				207,
				164,
				10,
				43,
				101
			],
			[
				65,
				237,
				112,
				85,
				29,
				215,
				232,
				65,
				136,
				58,
				184,
				240,
				177,
				107,
				240,
				65,
				118,
				183,
				209,
				72,
				14,
				79,
				10,
				249,
				243,
				212,
				195,
				89,
				87,
				104,
				208,
				104
			],
			[
				32,
				210,
				167,
				188,
				153,
				73,
				135,
				48,
				46,
				91,
				26,
				200,
				15,
				196,
				37,
				254,
				37,
				248,
				182,
				49,
				105,
				234,
				120,
				230,
				143,
				186,
				174,
				250,
				89,
				55,
				155,
				191
			]
		],
		"Flags": "tws="
	}
	`

	msg := bip37.ReadBlock(t)

	/*
		btcf := btcbloom.NewFilter(10, bloom.Tweak, 0.000001, btcwire.BloomUpdateAll)
		for i, tx := range msg.Transactions {
			if i%2 == 0 {
				continue
			}

			h := tx.TxHash()
			btcf.AddHash(&h)
		}

		btcBlock := btcutil.NewBlock(msg)
		B1, _ := btcbloom.NewMerkleBlock(btcBlock, btcf)

		expect, _ := json.MarshalIndent(B1, "", "  ")
		t.Logf("%s", expect)
	*/

	bf := bloom.New(10, 0.000001, wire.UpdateAll)
	for i, tx := range msg.Transactions {
		if i%2 == 0 {
			continue
		}

		h := tx.TxHash()
		bf.Add(h[:])
	}
	got, _ := merkle.New(msg, bf)

	expect := new(btcwire.MsgMerkleBlock)
	if err := json.Unmarshal([]byte(expectJSON), expect); nil != err {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expect) {
		t.Fatal("failed")
	}
}
