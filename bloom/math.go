package bloom

// MinUint32 estimates the larger one of x and y
func MinUint32(x, y uint32) uint32 {
	if x <= y {
		return x
	}

	return y
}
