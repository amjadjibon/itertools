package itertools

import (
	"iter"
	"sort"
)

// Iterator is a generic iterator that can be used
// to iterate over any type of sequence
type Iterator[V any] struct {
	seq  iter.Seq[V]
	curr *V
	done bool
}

// ToIter creates an Iterator from a slice
func ToIter[V any](slice []V) *Iterator[V] {
	ch := make(chan V)
	go func() {
		for _, v := range slice {
			ch <- v
		}
		close(ch)
	}()

	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for v := range ch {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Next advances the iterator and returns true if there is a next element.
func (it *Iterator[V]) Next() bool {
	if it.done {
		return false
	}

	var next V
	hasNext := false
	it.seq(func(v V) bool {
		next = v
		hasNext = true
		return false
	})

	if hasNext {
		it.curr = &next
		return true
	}

	it.done = true
	return false
}

// Current returns the current element of the iterator
func (it *Iterator[V]) Current() V {
	if it.curr == nil {
		panic("iterator is not started")
	}

	if it.done {
		panic("iterator is done")
	}

	return *it.curr
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

// Drop returns an Iterator that skips the first n elements of the Iterator
func (it *Iterator[V]) Drop(n int) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			i := 0
			it.seq(func(v V) bool {
				if i < n {
					i++
					return true
				}
				return yield(v)
			})
		},
	}
}

// TakeWhile returns an Iterator that yields elements while the predicate is true
func (it *Iterator[V]) TakeWhile(predicate func(V) bool) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			it.seq(func(v V) bool {
				if predicate(v) {
					return yield(v)
				}
				return false
			})
		},
	}
}

// DropWhile returns an Iterator that skips elements while the predicate is true
func (it *Iterator[V]) DropWhile(predicate func(V) bool) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			var dropping bool
			it.seq(func(v V) bool {
				if dropping {
					return yield(v)
				}
				if !predicate(v) {
					dropping = true
					return yield(v)
				}
				return true
			})
		},
	}
}

// First returns the first element of the Iterator
func (it *Iterator[V]) First() V {
	it.Next()
	return it.Current()
}

// Last returns the last element of the Iterator
func (it *Iterator[V]) Last() V {
	return it.Reverse().First()
}

// All returns true if all elements in the Iterator satisfy the predicate
func (it *Iterator[V]) All(predicate func(V) bool) bool {
	all := true
	it.seq(func(v V) bool {
		if !predicate(v) {
			all = false
			return false
		}
		return true
	})
	return all
}

// Any returns true if any element in the Iterator satisfies the predicate
func (it *Iterator[V]) Any(predicate func(V) bool) bool {
	any := false
	it.seq(func(v V) bool {
		if predicate(v) {
			any = true
			return false
		}
		return true
	})
	return any
}

// Find returns the first element that satisfies the predicate
func (it *Iterator[V]) Find(predicate func(V) bool) (V, bool) {
	var result V
	var found bool

	it.seq(func(v V) bool {
		if predicate(v) {
			result = v
			found = true
			return false
		}
		return true
	})

	return result, found
}

// Sort returns an Iterator with elements sorted in ascending order using the provided less function.
func (it *Iterator[V]) Sort(less func(a, b V) bool) *Iterator[V] {
	xs := it.Collect()
	sort.Slice(xs, func(i, j int) bool {
		return less(xs[i], xs[j])
	})
	return ToIter(xs)
}

// Min returns the minimum element in the Iterator using the provided less function
func (it *Iterator[V]) Min(less func(a, b V) bool) (V, bool) {
	var min V
	var found bool

	it.seq(func(v V) bool {
		if !found || less(v, min) {
			min = v
			found = true
		}
		return true
	})

	return min, found
}

// Max returns the maximum element in the Iterator using the provided less function
func (it *Iterator[V]) Max(less func(a, b V) bool) (V, bool) {
	var max V
	var found bool

	it.seq(func(v V) bool {
		if !found || less(max, v) {
			max = v
			found = true
		}
		return true
	})

	return max, found
}
