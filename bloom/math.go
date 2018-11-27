package bloom

func MinUint32(x, y uint32) uint32 {
	if x <= y {
		return x
	}

	return y
}
