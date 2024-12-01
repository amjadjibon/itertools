package itertools

import (
	"cmp"

	"golang.org/x/exp/constraints"
)

// Zip combines two iterators element-wise into a single iterator of pairs.
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
		}) bool) {
			ch1 := make(chan A)
			ch2 := make(chan B)
			go func() { it1.seq(func(v A) bool { ch1 <- v; return true }); close(ch1) }()
			go func() { it2.seq(func(v B) bool { ch2 <- v; return true }); close(ch2) }()
			for v1 := range ch1 {
				v2, ok := <-ch2
				if !ok {
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
// If one iterator is longer than the other, the shorter iterator is extended with the fill value.
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
		}) bool) {
			ch1 := make(chan A)
			ch2 := make(chan B)
			go func() { it1.seq(func(v A) bool { ch1 <- v; return true }); close(ch1) }()
			go func() { it2.seq(func(v B) bool { ch2 <- v; return true }); close(ch2) }()
			for {
				v1, ok1 := <-ch1
				v2, ok2 := <-ch2
				if !ok1 && !ok2 {
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

// Fold accumulates the elements of the iterator
func Fold[V any, T any](it *Iterator[V], transform func(T, V) T, initial T) T {
	acc := initial

	it.seq(func(v V) bool {
		acc = transform(acc, v)
		return true
	})

	return acc
}

// Sum adds all elements of the iterator
func Sum[V any, T cmp.Ordered](it *Iterator[V], transform func(V) T, zero T) T {
	return Fold(it, func(acc T, v V) T { return acc + transform(v) }, zero)
}

type Productable interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Product multiplies all elements of the iterator
func Product[V any, T Productable](it *Iterator[V], transform func(V) T, one T) T {
	return Fold(it, func(acc T, v V) T { return acc * transform(v) }, one)
}
