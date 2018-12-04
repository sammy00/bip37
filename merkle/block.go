package merkle

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/sammy00/bip37/bloom"
)

type Block struct{}

// the width of tree as height is estimated as
//  ceil(#(leaves)/2^k)=(#(leaves)+2^k-1)/2^k
func (block *Block) calcTreeWidth(height uint32) uint32 {
	return 0
}

func (block *Block) hash(height, idx uint32) *chainhash.Hash {
	return nil
}

// traverse and build the depth-first sub-tree of the given height and
// indexed by idx within that row, where the index of first node of each
// row is 0
func (block *Block) traverseAndBuild(height, idx uint32) {}

func New(block *wire.MsgBlock, filter *bloom.Filter) (*wire.MsgMerkleBlock,
	[]uint32) {
	// retrieve all txs
	// calculates digests for all leaf txs
	// filter out the matched txs

	// calculate the tree height

	// build the depth-first partial Merkle tree

	// convert the native block to the canonical one, which would
	//  + add all tx hashes
	//  + populate the flag bits

	return nil, nil
}

func Validate(block *wire.MsgMerkleBlock) bool {
	return false
}
