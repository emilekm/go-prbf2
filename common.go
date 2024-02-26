package prdemo

import (
	"cmp"
	"io"
	"sort"

	"github.com/ghostiam/binstruct"
)

func newBinReader(r io.ReadSeeker) binstruct.Reader {
	return binstruct.NewReader(r, demoEndian, false)
}

func sortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
