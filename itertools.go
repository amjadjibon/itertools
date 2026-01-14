package itertools

import (
	"cmp"
	"iter"

	"golang.org/x/exp/constraints"
)

// NewIterator creates a new iterator from a variadic list of values.
//
// Example:
//
//	iter := itertools.NewIterator(1, 2, 3, 4, 5)
//	result := iter.Collect()
//	// result is []int{1, 2, 3, 4, 5}
func NewIterator[V any](v ...V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for _, v := range v {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Zip combines two iterators element-wise into a single iterator of pairs.
// Iteration stops when either iterator is exhausted.
//
// The resulting iterator yields struct values with First and Second fields.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3})
//	iter2 := itertools.ToIter([]string{"a", "b", "c"})
//	zipped := itertools.Zip(iter1, iter2).Collect()
//	// zipped is [{First: 1, Second: "a"}, {First: 2, Second: "b"}, {First: 3, Second: "c"}]
func Zip[A, B any](it1 *Iterator[A], it2 *Iterator[B]) *Iterator[struct {
	First  A
	Second B
}] {
	return &Iterator[struct {
		First  A
		Second B
	}]{
		seq: func(yield func(struct {
			First  A
			Second B
		}) bool,
		) {
			pull1, stop1 := iter.Pull(it1.seq)
			pull2, stop2 := iter.Pull(it2.seq)
			defer stop1()
			defer stop2()

			for {
				v1, ok1 := pull1()
				v2, ok2 := pull2()
				if !ok1 || !ok2 {
					return
				}
				if !yield(struct {
					First  A
					Second B
				}{v1, v2}) {
					return
				}
			}
		},
	}
}

// Zip2 combines two iterators element-wise into a single iterator of pairs.
// Unlike Zip, if one iterator is longer than the other, the shorter one is
// extended using the fill values.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3})
//	iter2 := itertools.ToIter([]string{"a"})
//	fill := struct{ First int; Second string }{0, ""}
//	zipped := itertools.Zip2(iter1, iter2, fill).Collect()
//	// zipped is [{1, "a"}, {2, ""}, {3, ""}]
func Zip2[A, B any](it1 *Iterator[A], it2 *Iterator[B], fill struct {
	First  A
	Second B
}) *Iterator[struct {
	First  A
	Second B
}] {
	return &Iterator[struct {
		First  A
		Second B
	}]{
		seq: func(yield func(struct {
			First  A
			Second B
		}) bool,
		) {
			pull1, stop1 := iter.Pull(it1.seq)
			pull2, stop2 := iter.Pull(it2.seq)
			defer stop1()
			defer stop2()

			for {
				v1, ok1 := pull1()
				v2, ok2 := pull2()

				if !ok1 && !ok2 {
					return
				}

				// Use fill values when one iterator ends
				result := struct {
					First  A
					Second B
				}{}

				if ok1 {
					result.First = v1
				} else {
					result.First = fill.First
				}

				if ok2 {
					result.Second = v2
				} else {
					result.Second = fill.Second
				}

				if !yield(result) {
					return
				}
			}
		},
	}
}

// Fold accumulates the elements of the iterator using a binary operation.
// Also known as reduce or aggregate in other languages.
//
// The transform function takes an accumulator and the next value,
// returning the new accumulator value.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	sum := itertools.Fold(iter, func(acc, v int) int { return acc + v }, 0)
//	// sum is 15
func Fold[V any, T any](it *Iterator[V], transform func(T, V) T, initial T) T {
	acc := initial

	it.seq(func(v V) bool {
		acc = transform(acc, v)
		return true
	})

	return acc
}

// Sum adds all elements of the iterator after applying the transform function.
// The zero parameter specifies the additive identity for the result type.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	sum := itertools.Sum(iter, func(v int) int { return v }, 0)
//	// sum is 15
//
//	// Sum of squares
//	iter2 := itertools.ToIter([]int{1, 2, 3})
//	sumSquares := itertools.Sum(iter2, func(v int) int { return v * v }, 0)
//	// sumSquares is 14
func Sum[V any, T cmp.Ordered](it *Iterator[V], transform func(V) T, zero T) T {
	return Fold(it, func(acc T, v V) T { return acc + transform(v) }, zero)
}

// Productable is a constraint for types that support multiplication.
type Productable interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Product multiplies all elements of the iterator after applying the transform function.
// The one parameter specifies the multiplicative identity for the result type.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	product := itertools.Product(iter, func(v int) int { return v }, 1)
//	// product is 120
func Product[V any, T Productable](it *Iterator[V], transform func(V) T, one T) T {
	return Fold(it, func(acc T, v V) T { return acc * transform(v) }, one)
}

// ChunkSlice returns an iterator that yields slices of up to `size` elements.
// The last chunk may contain fewer elements if the total is not divisible by size.
//
// Each chunk is a separate slice - modifications won't affect the original data.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5, 6, 7})
//	chunks := itertools.ChunkSlice(iter, 3).Collect()
//	// chunks is [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
func ChunkSlice[V any](it *Iterator[V], size int) *Iterator[[]V] {
	return &Iterator[[]V]{
		seq: func(yield func([]V) bool) {
			chunk := make([]V, 0, size)
			it.seq(func(v V) bool {
				chunk = append(chunk, v)
				if len(chunk) == size {
					// Create a new slice to avoid sharing the underlying array
					result := make([]V, size)
					copy(result, chunk)
					if !yield(result) {
						return false
					}
					chunk = chunk[:0] // Clear the chunk while preserving capacity
				}
				return true
			})
			// Handle any remaining elements in the last chunk
			if len(chunk) > 0 {
				result := make([]V, len(chunk))
				copy(result, chunk)
				yield(result)
			}
		},
	}
}

// Chunks returns an iterator that yields iterators, each containing up to `size` elements.
// The last chunk may contain fewer elements if the total is not divisible by size.
//
// Unlike ChunkSlice, this returns iterators instead of slices.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	chunks := itertools.Chunks(iter, 2)
//	for chunks.Next() {
//	    chunk := chunks.Current()
//	    fmt.Println(chunk.Collect())
//	}
//	// Output: [1 2]
//	//         [3 4]
//	//         [5]
func Chunks[V any](it *Iterator[V], size int) *Iterator[*Iterator[V]] {
	return &Iterator[*Iterator[V]]{
		seq: func(yield func(*Iterator[V]) bool) {
			chunk := make([]V, 0, size)
			it.seq(func(v V) bool {
				chunk = append(chunk, v)
				if len(chunk) == size {
					// Create a new slice to avoid sharing the underlying array
					result := make([]V, size)
					copy(result, chunk)
					if !yield(ToIter(result)) {
						return false
					}
					chunk = chunk[:0] // Clear the chunk while preserving capacity
				}
				return true
			})
			// Handle any remaining elements in the last chunk
			if len(chunk) > 0 {
				result := make([]V, len(chunk))
				copy(result, chunk)
				yield(ToIter(result))
			}
		},
	}
}

// ChunkList returns a slice of iterators, each containing up to `size` elements.
// This is a convenience function that collects Chunks into a slice.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5, 6})
//	chunks := itertools.ChunkList(iter, 2)
//	// chunks is []*Iterator with 3 iterators containing [1,2], [3,4], [5,6]
func ChunkList[V any](it *Iterator[V], size int) []*Iterator[V] {
	return Chunks(it, size).Collect()
}

// Flatten concatenates multiple iterators into a single iterator.
// Elements are yielded in order: all elements from the first iterator,
// then all from the second, and so on.
//
// Properly handles early termination - stops immediately when yield returns false.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3})
//	iter2 := itertools.ToIter([]int{4, 5, 6})
//	iter3 := itertools.ToIter([]int{7, 8, 9})
//	flattened := itertools.Flatten(iter1, iter2, iter3).Collect()
//	// flattened is []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
func Flatten[V any](its ...*Iterator[V]) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for _, it := range its {
				shouldContinue := true
				it.seq(func(v V) bool {
					if !yield(v) {
						shouldContinue = false
						return false
					}
					return true
				})
				if !shouldContinue {
					return
				}
			}
		},
	}
}

// CartesianProduct returns an iterator of all pairs of elements from two iterators.
// Note: The second iterator (it2) is fully collected into memory to enable
// multiple iterations over it for each element in it1. Use with caution for large datasets.
//
// Example:
//
//	iter1 := ToIter([]int{1, 2})
//	iter2 := ToIter([]string{"a", "b"})
//	product := CartesianProduct(iter1, iter2).Collect()
//	// Result: [{1, "a"}, {1, "b"}, {2, "a"}, {2, "b"}]
func CartesianProduct[A, B any](it1 *Iterator[A], it2 *Iterator[B]) *Iterator[struct {
	X A
	Y B
}] {
	xs := it2.Collect()
	return &Iterator[struct {
		X A
		Y B
	}]{
		seq: func(yield func(struct {
			X A
			Y B
		}) bool,
		) {
			it1.seq(func(a A) bool {
				for _, b := range xs {
					if !yield(struct {
						X A
						Y B
					}{X: a, Y: b}) {
						return false
					}
				}
				return true
			})
		},
	}
}
