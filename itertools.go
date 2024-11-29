package itertools

import (
	"iter"
)

// Iterator is a generic iterator that can be used
// to iterate over any type of sequence
type Iterator[V any] struct {
	seq iter.Seq[V]
}

// ToIter creates an Iterator from a slice
func ToIter[V any](slice []V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Collect collects all elements from the Iterator into a slice.
func (it *Iterator[V]) Collect() []V {
	collect := make([]V, 0)
	it.seq(func(e V) bool {
		collect = append(collect, e)
		return true
	})
	return collect
}

// Each applies a function to each element of the Iterator.
func (it *Iterator[V]) Each(f func(V)) {
	it.seq(func(v V) bool {
		f(v)
		return true
	})
}
