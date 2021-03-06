package merkle

import (
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/sammyne/bip37/bloom"
)

// Block is an intermediate data structure helping to build a merkle block
// based on a given bloom filter
type Block struct {
	*btcutil.Block

	flags  []byte
	leaves []*chainhash.Hash
	// included signals where a merkle leaf possibly matches the filter
	included []byte
	branches []*chainhash.Hash
	nTx      uint32
}

// calcTreeWidth estimates the width of tree at height according to
//  width = ceil(#(leaves)/2^h)=(#(leaves)+2^h-1)/2^h
// where the height of leaves is defined as 0
func (block *Block) calcTreeWidth(height uint32) uint32 {
	return (block.nTx + (1 << height) - 1) >> height
}

// branchHash calculates the hash of branch by post-order traversing rooted at
// the idx-th (zero-based) node of height, where
//  - the hash for a leaf is defined to be itself
//  - for the hash of a node has no right branch is defines as H(L|L)
// **NOTE** the initial provided index is assumed to be within bound
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

// traverseAndBuild traverses and builds the depth-first sub-tree of the given
// height and indexed by idx within that row, where the index of first
// node of each row is 0.
// The defailed algorithm is specified as https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#constructing-a-partial-merkle-tree-object.
func (block *Block) traverseAndBuild(height, idx uint32) {
	var flag byte
	// For a subtree rooted at the idx-th node of height
	// 	- (idx<<height) is the index of the leftmost leaf
	//  - all leaves of index (i>>height==idx) is the leaves of this subtree
	// The following loop is to find if any leaf of this subtree matching the
	// filter, thus assigning the proper flag to this node.
	// In case of no matching leaf, the flag would just remain 0
	for i := idx << height; (i < block.nTx) && (i>>height == idx) &&
		(0x01 != flag); i++ {
		flag = block.included[i]
	}

	// append the flag for this node
	block.flags = append(block.flags, flag)

	// The sibling leaves of included leaves must be included.
	// The hash for parent of non-included leaves must be included.
	// These 2 cases is base cases
	if 0 == height || 0x00 == flag {
		block.branches = append(block.branches, block.branchHash(height, idx))
		return
	}

	// left child branch
	block.traverseAndBuild(height-1, idx<<1)
	// right child branch
	if j := (idx << 1) + 1; j < block.calcTreeWidth(height-1) {
		block.traverseAndBuild(height-1, j)
	}
}

// New builds a merkle block based on the given raw block and filter
func New(b *wire.MsgBlock, filter *bloom.Filter) (*wire.MsgMerkleBlock,
	[]uint32) {
	nTx := len(b.Transactions)
	block := &Block{
		Block: btcutil.NewBlock(b),

		included: make([]byte, nTx),
		leaves:   make([]*chainhash.Hash, nTx),
		nTx:      uint32(nTx),
	}

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
		Hashes:       make([]*chainhash.Hash, 0, len(block.branches)),
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

// Parse validates if the given merkle block is valid, and returns the
// matching hash if any
func Parse(block *wire.MsgMerkleBlock) ([]*chainhash.Hash, bool) {
	// calculate the tree height
	var height uint32
	for ; (1 << height) < block.Transactions; height++ {
	}

	var (
		j, k    int
		matched []*chainhash.Hash
	)

	root := parse(&matched, block, 0, height, &j, &k)

	// Check
	//  - all hashes have been consumed
	//  - all flag bits have been consumed
	//  - the merkle root matches
	ok := len(block.Hashes) == k &&
		len(block.Flags) == (j+7)/8 &&
		0 == (block.Flags[j>>3]>>uint(j%8)) &&
		block.Header.MerkleRoot.IsEqual(root)

		//fmt.Println(len(block.Hashes) == k)
	//fmt.Println(len(block.Flags) == (j+7)/8)
	//fmt.Println(0 == (block.Flags[j>>3] >> uint(j%8)))
	// check the PoW

	return matched, ok
}

// calcTreeWidth calculates the number of nodes at height for a tree
// with nTx leaves
func calcTreeWidth(nTx, height uint32) uint32 {
	return (nTx + (1 << height) - 1) >> height
}

// j is the #(flag-bit) consumed
// k is the #(hash) consumed
func parse(matched *[]*chainhash.Hash, block *wire.MsgMerkleBlock,
	i, height uint32, j, k *int) *chainhash.Hash {
	if (*j>>3) >= len(block.Flags) || *k >= len(block.Hashes) {
		// flag bits or hash list is exhausted
		return nil
	}

	flag := (block.Flags[*j>>3] >> uint(*j%8)) & 0x01
	*j++
	if 0 == flag { // the non-included nodes
		hash := block.Hashes[*k]
		*k++

		return hash
	} else if 0 == height { // a included leaf
		hash := block.Hashes[*k]
		*k++
		*matched = append(*matched, hash)

		return hash
	}

	childIdx := i << 1
	L := parse(matched, block, childIdx, height-1, j, k)
	if nil == L {
		return nil
	}

	childIdx++
	// the missing right branch is replaced with its left sibling
	if childIdx >= calcTreeWidth(block.Transactions, height-1) {
		return blockchain.HashMerkleBranches(L, L)
	}

	R := parse(matched, block, childIdx, height-1, j, k)
	if nil != R && !R.IsEqual(L) {
		return blockchain.HashMerkleBranches(L, R)
	}

	return nil
}
