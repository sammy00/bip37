package bloom_test

import (
	"testing"

	"github.com/sammy00/bip37/bloom"
)

func TestMinUint32(t *testing.T) {
	testCases := []struct {
		x, y   uint32
		expect uint32
	}{
		{1, 2, 1},
		{3, 2, 2},
		{2, 2, 2},
	}

	for i, c := range testCases {
		if got := bloom.MinUint32(c.x, c.y); got != c.expect {
			t.Fatalf("#%d failed: got %d, expect %d", i, got, c.expect)
		}
	}
}
