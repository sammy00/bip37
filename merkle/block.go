package merkle

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/sammy00/bip37/bloom"
)

type Block struct {
	*btcutil.Block

	flags    []byte
	leaves   []*chainhash.Hash
	included []byte
	branches []*chainhash.Hash
	nTx      uint32
}

// the width of tree as height is estimated as
//  ceil(#(leaves)/2^h)=(#(leaves)+2^h-1)/2^h
func (block *Block) calcTreeWidth(height uint32) uint32 {
	return (block.nTx + (1 << height) - 1) >> height
}
func (block *Block) branchHash(height, idx uint32) *chainhash.Hash {
	if 0 == height {
		return block.leaves[idx]
	}

	var L, R *chainhash.Hash
	L = block.branchHash(height-1, idx<<1)
	if j := (idx << 1) + 1; j < block.calcTreeWidth(height-1) {
		R = block.branchHash(height-1, j)
	} else {
		R = L
	}

	return blockchain.HashMerkleBranches(L, R)
}

// traverse and build the depth-first sub-tree of the given height and
// indexed by idx within that row, where the index of first node of each
// row is 0
func (block *Block) traverseAndBuild(height, idx uint32) {
	var flag byte
	for i := idx << height; (i < block.nTx) && (i>>height == idx) &&
		(0x01 != flag); i++ {
		flag = block.included[i]
	}

	block.flags = append(block.flags, flag)

	if 0 == height || 0x00 == flag {
		block.branches = append(block.branches, block.branchHash(height, idx))
		return
	}

	block.traverseAndBuild(height-1, idx<<1)
	if j := (idx << 1) + 1; j < block.calcTreeWidth(height-1) {
		block.traverseAndBuild(height-1, j)
	}
}

func New(b *wire.MsgBlock, filter *bloom.Filter) (*wire.MsgMerkleBlock,
	[]uint32) {
	block := &Block{Block: btcutil.NewBlock(b)}

	// retrieve all txs
	block.nTx = uint32(len(block.Transactions()))
	block.included = make([]byte, block.nTx)

	var hits []uint32
	// calculates digests for all leaf txs
	for i, tx := range block.Transactions() {
		block.leaves[i] = tx.Hash()
		if filter.MatchTxAndUpdate(tx) {
			// filter out the matched txs and append matched bit
			block.included[i] = 0x01
			hits = append(hits, uint32(i))
		} else {
			block.included[i] = 0x00
		}
	}

	// calculate the tree height
	var height uint32
	for ; (1 << height) < block.nTx; height++ {
	}

	// build the depth-first partial Merkle tree
	block.traverseAndBuild(height, 0)

	// convert the native block to the canonical one, which would
	//  + add all tx hashes
	//  + populate the flag bits
	msg := &wire.MsgMerkleBlock{
		Hashes:       make([]*chainhash.Hash, len(block.branches)),
		Header:       block.MsgBlock().Header,
		Transactions: block.nTx,
		Flags:        make([]byte, (len(block.flags)+7)/8),
	}
	for _, m := range block.branches {
		msg.AddTxHash(m)
	}
	for i, b := range block.flags {
		msg.Flags[i/8] |= b << uint32(i%8)
	}

	return msg, hits
}

func Validate(block *wire.MsgMerkleBlock) bool {
	return false
}
