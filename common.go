package prdemo

import (
	"cmp"
	"io"
	"sort"
	"time"

	"github.com/ghostiam/binstruct"
	"golang.org/x/exp/constraints"
)

type Time[T constraints.Unsigned] struct {
	time.Time
}

func (t *Time[T]) Decode(m *Message) error {
	var ts T
	err := m.Decode(&ts)
	if err != nil {
		return err
	}

	t = &Time[T]{time.Unix(int64(ts), 0)}

	return nil
}

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
