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

// Reverse returns an Iterator that iterates over the elements in reverse order
func (it *Iterator[V]) Reverse() *Iterator[V] {
	xs := it.Collect()
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
	return ToIter(xs)
}

// Filter returns an Iterator that only yields elements that satisfy the predicate
func (it *Iterator[V]) Filter(predicate func(V) bool) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			it.seq(func(v V) bool {
				if predicate(v) {
					return yield(v)
				}
				return true
			})
		},
	}
}

// Map transforms each element in the Iterator using a provided function.
func (it *Iterator[V]) Map(f func(V) V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			it.seq(func(v V) bool {
				return yield(f(v))
			})
		},
	}
}

// Chain concatenates two iterators
func (it *Iterator[V]) Chain(other *Iterator[V]) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			it.seq(yield)
			other.seq(yield)
		},
	}
}

// Take returns an Iterator that yields the first n elements of the Iterator
func (it *Iterator[V]) Take(n int) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			i := 0
			it.seq(func(v V) bool {
				if i < n {
					i++
					return yield(v)
				}
				return false
			})
		},
	}
}
