package wire

// FilterAdd is the command to add more patterns to the filter
// Detail sees https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#new-messages
type FilterAdd struct {
	Data []byte
}

// FilterLoad defines the command sent to the full peer node to initialize a
// connection capable of filtering tx.
// Detail sees https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#new-messages
type FilterLoad struct {
	// Bits records the matching pattern as a bit field
	Bits []byte
	// HashFuncs is the number of hash functions to use
	HashFuncs uint32
	// Tweak is a random value add to the seed value in the hash function
	Tweak uint32
	// Flags specifies the updating policy of the filter
	Flags BloomUpdateType
}

// FilterClear defines the command to reset the filter.
// Detail sees https://github.com/bitcoin/bips/blob/master/bip-0037.mediawiki#new-messages
type FilterClear struct{}
