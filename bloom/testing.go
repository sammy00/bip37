package bloom

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func FakeBlock(t *testing.T) *btcutil.Block {
	blockStr := "0100000082bb869cf3a793432a66e826e05a6fc37469f8efb7421dc" +
		"880670100000000007f16c5962e8bd963659c793ce370d95f093bc7e367" +
		"117b3c30c1f8fdd0d9728776381b4d4c86041b554b85290701000000010" +
		"00000000000000000000000000000000000000000000000000000000000" +
		"0000ffffffff07044c86041b0136ffffffff0100f2052a0100000043410" +
		"4eaafc2314def4ca98ac970241bcab022b9c1e1f4ea423a20f134c876f2" +
		"c01ec0f0dd5b2e86e7168cefe0d81113c3807420ce13ad1357231a22522" +
		"47d97a46a91ac000000000100000001bcad20a6a29827d1424f08989255" +
		"120bf7f3e9e3cdaaa6bb31b0737fe048724300000000494830450220356" +
		"e834b046cadc0f8ebb5a8a017b02de59c86305403dad52cd77b55af062e" +
		"a10221009253cd6c119d4729b77c978e1e2aa19f5ea6e0e52b3f16e32fa" +
		"608cd5bab753901ffffffff02008d380c010000001976a9142b4b8072ec" +
		"bba129b6453c63e129e643207249ca88ac0065cd1d000000001976a9141" +
		"b8dd13b994bcfc787b32aeadf58ccb3615cbd5488ac0000000001000000" +
		"03fdacf9b3eb077412e7a968d2e4f11b9a9dee312d666187ed77ee7d26a" +
		"f16cb0b000000008c493046022100ea1608e70911ca0de5af51ba57ad23" +
		"b9a51db8d28f82c53563c56a05c20f5a87022100a8bdc8b4a8acc8634c6" +
		"b420410150775eb7f2474f5615f7fccd65af30f310fbf01410465fdf49e" +
		"29b06b9a1582287b6279014f834edc317695d125ef623c1cc3aaece245b" +
		"d69fcad7508666e9c74a49dc9056d5fc14338ef38118dc4afae5fe2c585" +
		"caffffffff309e1913634ecb50f3c4f83e96e70b2df071b497b8973a3e7" +
		"5429df397b5af83000000004948304502202bdb79c596a9ffc24e96f438" +
		"6199aba386e9bc7b6071516e2b51dda942b3a1ed022100c53a857e76b72" +
		"4fc14d45311eac5019650d415c3abb5428f3aae16d8e69bec2301ffffff" +
		"ff2089e33491695080c9edc18a428f7d834db5b6d372df13ce2b1b0e0cb" +
		"cb1e6c10000000049483045022100d4ce67c5896ee251c810ac1ff9cecc" +
		"d328b497c8f553ab6e08431e7d40bad6b5022033119c0c2b7d792d31f11" +
		"87779c7bd95aefd93d90a715586d73801d9b47471c601ffffffff010071" +
		"4460030000001976a914c7b55141d097ea5df7a0ed330cf794376e53ec8" +
		"d88ac0000000001000000045bf0e214aa4069a3e792ecee1e1bf0c1d397" +
		"cde8dd08138f4b72a00681743447000000008b48304502200c45de8c4f3" +
		"e2c1821f2fc878cba97b1e6f8807d94930713aa1c86a67b9bf1e4022100" +
		"8581abfef2e30f957815fc89978423746b2086375ca8ecf359c85c2a5b7" +
		"c88ad01410462bb73f76ca0994fcb8b4271e6fb7561f5c0f9ca0cf64852" +
		"61c4a0dc894f4ab844c6cdfb97cd0b60ffb5018ffd6238f4d87270efb1d" +
		"3ae37079b794a92d7ec95ffffffffd669f7d7958d40fc59d2253d88e0f2" +
		"48e29b599c80bbcec344a83dda5f9aa72c000000008a473044022078124" +
		"c8beeaa825f9e0b30bff96e564dd859432f2d0cb3b72d3d5d93d38d7e93" +
		"0220691d233b6c0f995be5acb03d70a7f7a65b6bc9bdd426260f38a1346" +
		"669507a3601410462bb73f76ca0994fcb8b4271e6fb7561f5c0f9ca0cf6" +
		"485261c4a0dc894f4ab844c6cdfb97cd0b60ffb5018ffd6238f4d87270e" +
		"fb1d3ae37079b794a92d7ec95fffffffff878af0d93f5229a68166cf051" +
		"fd372bb7a537232946e0a46f53636b4dafdaa4000000008c49304602210" +
		"0c717d1714551663f69c3c5759bdbb3a0fcd3fab023abc0e522fe6440de" +
		"35d8290221008d9cbe25bffc44af2b18e81c58eb37293fd7fe1c2e7b46f" +
		"c37ee8c96c50ab1e201410462bb73f76ca0994fcb8b4271e6fb7561f5c0" +
		"f9ca0cf6485261c4a0dc894f4ab844c6cdfb97cd0b60ffb5018ffd6238f" +
		"4d87270efb1d3ae37079b794a92d7ec95ffffffff27f2b668859cd7f2f8" +
		"94aa0fd2d9e60963bcd07c88973f425f999b8cbfd7a1e2000000008c493" +
		"046022100e00847147cbf517bcc2f502f3ddc6d284358d102ed20d47a8a" +
		"a788a62f0db780022100d17b2d6fa84dcaf1c95d88d7e7c30385aecf415" +
		"588d749afd3ec81f6022cecd701410462bb73f76ca0994fcb8b4271e6fb" +
		"7561f5c0f9ca0cf6485261c4a0dc894f4ab844c6cdfb97cd0b60ffb5018" +
		"ffd6238f4d87270efb1d3ae37079b794a92d7ec95ffffffff0100c817a8" +
		"040000001976a914b6efd80d99179f4f4ff6f4dd0a007d018c385d2188a" +
		"c000000000100000001834537b2f1ce8ef9373a258e10545ce5a50b758d" +
		"f616cd4356e0032554ebd3c4000000008b483045022100e68f422dd7c34" +
		"fdce11eeb4509ddae38201773dd62f284e8aa9d96f85099d0b002202243" +
		"bd399ff96b649a0fad05fa759d6a882f0af8c90cf7632c2840c29070aec" +
		"20141045e58067e815c2f464c6a2a15f987758374203895710c2d452442" +
		"e28496ff38ba8f5fd901dc20e29e88477167fe4fc299bf818fd0d9e1632" +
		"d467b2a3d9503b1aaffffffff0280d7e636030000001976a914f34c3e10" +
		"eb387efe872acb614c89e78bfca7815d88ac404b4c00000000001976a91" +
		"4a84e272933aaf87e1715d7786c51dfaeb5b65a6f88ac00000000010000" +
		"000143ac81c8e6f6ef307dfe17f3d906d999e23e0189fda838c5510d850" +
		"927e03ae7000000008c4930460221009c87c344760a64cb8ae6685a3eec" +
		"2c1ac1bed5b88c87de51acd0e124f266c16602210082d07c037359c3a25" +
		"7b5c63ebd90f5a5edf97b2ac1c434b08ca998839f346dd40141040ba7e5" +
		"21fa7946d12edbb1d1e95a15c34bd4398195e86433c92b431cd315f455f" +
		"e30032ede69cad9d1e1ed6c3c4ec0dbfced53438c625462afb792dcb098" +
		"544bffffffff0240420f00000000001976a9144676d1b820d63ec272f19" +
		"00d59d43bc6463d96f888ac40420f00000000001976a914648d04341d00" +
		"d7968b3405c034adc38d4d8fb9bd88ac00000000010000000248cc91750" +
		"1ea5c55f4a8d2009c0567c40cfe037c2e71af017d0a452ff705e3f10000" +
		"00008b483045022100bf5fdc86dc5f08a5d5c8e43a8c9d5b1ed8c65562e" +
		"280007b52b133021acd9acc02205e325d613e555f772802bf413d36ba80" +
		"7892ed1a690a77811d3033b3de226e0a01410429fa713b124484cb2bd7b" +
		"5557b2c0b9df7b2b1fee61825eadc5ae6c37a9920d38bfccdc7dc3cb0c4" +
		"7d7b173dbc9db8d37db0a33ae487982c59c6f8606e9d1791ffffffff41e" +
		"d70551dd7e841883ab8f0b16bf04176b7d1480e4f0af9f3d4c3595768d0" +
		"68000000008b4830450221008513ad65187b903aed1102d1d0c47688127" +
		"658c51106753fed0151ce9c16b80902201432b9ebcb87bd04ceb2de6603" +
		"5fbbaf4bf8b00d1cfe41f1a1f7338f9ad79d210141049d4cf80125bf50b" +
		"e1709f718c07ad15d0fc612b7da1f5570dddc35f2a352f0f27c978b0682" +
		"0edca9ef982c35fda2d255afba340068c5035552368bc7200c1488fffff" +
		"fff0100093d00000000001976a9148edb68822f1ad580b043c7b3df2e40" +
		"0f8699eb4888ac00000000"
	blockBytes, err := hex.DecodeString(blockStr)
	if err != nil {
		t.Fatalf("TestFilterInsertP2PubKeyOnly DecodeString failed: %v", err)
	}
	block, err := btcutil.NewBlockFromBytes(blockBytes)
	if err != nil {
		t.Fatalf("TestFilterInsertP2PubKeyOnly NewBlockFromBytes failed: %v", err)
	}

	out, _ := json.MarshalIndent(block.MsgBlock(), "", "  ")
	fmt.Printf("%s\n", out)

	return block
}

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

func ReadBlock(t *testing.T) *btcwire.MsgBlock {
	fd, err := os.Open(filepath.Join("testdata", "block.json"))
	if nil != err {
		t.Fatal(err)
	}
	defer fd.Close()

	msgBlock := new(btcwire.MsgBlock)
	if err := json.NewDecoder(fd).Decode(msgBlock); nil != err {
		t.Fatal(err)
	}

	return msgBlock
}
