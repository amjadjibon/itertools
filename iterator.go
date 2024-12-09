package itertools

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"reflect"
	"sort"
	"strings"
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

// Nth returns the nth element of the Iterator
func (it *Iterator[V]) Nth(n int) V {
	return it.Drop(n).First()
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
	anyFlag := false
	it.seq(func(v V) bool {
		if predicate(v) {
			anyFlag = true
			return false
		}
		return true
	})
	return anyFlag
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

// Partition returns two Iterators, one with elements that satisfy the predicate and one with elements that don't
func (it *Iterator[V]) Partition(predicate func(V) bool) (matched *Iterator[V], unmatched *Iterator[V]) {
	var yes []V
	var no []V

	it.seq(func(v V) bool {
		if predicate(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
		return true
	})

	matched = ToIter(yes)
	unmatched = ToIter(no)
	return
}

// Count returns the number of elements in the Iterator
func (it *Iterator[V]) Count() int {
	var count int
	it.seq(func(V) bool {
		count++
		return true
	})
	return count
}

// Unique returns an Iterator with only unique elements
func (it *Iterator[V]) Unique(keyFunc func(V) any) *Iterator[V] {
	seen := make(map[any]struct{})
	return it.Filter(func(v V) bool {
		key := keyFunc(v)
		if _, ok := seen[key]; ok {
			return false
		}
		seen[key] = struct{}{}
		return true
	})
}

// GroupBy groups elements by a key function into a map.
func (it *Iterator[V]) GroupBy(keyFunc func(V) string) map[string][]V {
	groups := make(map[string][]V)
	it.seq(func(v V) bool {
		key := keyFunc(v)
		groups[key] = append(groups[key], v)
		return true
	})
	return groups
}

// Cycle returns an Iterator that cycles through the elements of the Iterator indefinitely
func (it *Iterator[V]) Cycle() *Iterator[V] {
	xs := it.Collect()
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for {
				for _, v := range xs {
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// Repeat returns an Iterator that yields the same element n times
func Repeat[V any](v V, n int) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for i := 0; i < n; i++ {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Union returns an Iterator that yields elements from both iterators without duplicates
func (it *Iterator[V]) Union(other *Iterator[V], keyFunc func(V) any) *Iterator[V] {
	seen := make(map[any]struct{})
	return it.Chain(other).Filter(func(v V) bool {
		key := keyFunc(v)
		if _, ok := seen[key]; ok {
			return false
		}
		seen[key] = struct{}{}
		return true
	})
}

// Intersection returns an Iterator that yields elements that are present in both iterators
func (it *Iterator[V]) Intersection(other *Iterator[V], keyFunc func(V) any) *Iterator[V] {
	seen := make(map[any]struct{})
	other.seq(func(v V) bool {
		seen[keyFunc(v)] = struct{}{}
		return true
	})

	return it.Filter(func(v V) bool {
		_, ok := seen[keyFunc(v)]
		return ok
	})
}

// Difference returns an Iterator that yields elements that are present in the first iterator but not in the second
func (it *Iterator[V]) Difference(other *Iterator[V], keyFunc func(V) any) *Iterator[V] {
	seen := make(map[any]struct{})
	other.seq(func(v V) bool {
		seen[keyFunc(v)] = struct{}{}
		return true
	})

	return it.Filter(func(v V) bool {
		_, ok := seen[keyFunc(v)]
		return !ok
	})
}

// StepBy returns an Iterator that yields every nth element
func (it *Iterator[V]) StepBy(n int) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			i := 0
			it.seq(func(v V) bool {
				if i%n == 0 {
					if !yield(v) {
						return false
					}
				}
				i++
				return true
			})
		},
	}
}

// Shuffle returns an Iterator that yields elements in a random order
func (it *Iterator[V]) Shuffle() *Iterator[V] {
	xs := it.Collect()
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for _, i := range rand.Perm(len(xs)) {
				if !yield(xs[i]) {
					return
				}
			}
		},
	}
}

// Index returns the index of the first element that satisfies the predicate
func (it *Iterator[V]) Index(predicate func(V) bool) int {
	index := 0
	found := false
	it.seq(func(v V) bool {
		if predicate(v) {
			found = true
			return false
		}
		index++
		return true
	})
	if found {
		return index
	}
	return -1
}

// LastIndex returns the index of the last element that satisfies the predicate
func (it *Iterator[V]) LastIndex(predicate func(V) bool) int {
	index := -1
	i := 0
	it.seq(func(v V) bool {
		if predicate(v) {
			index = i
		}
		i++
		return true
	})
	return index
}

// IsSorted returns true if the elements in the Iterator are sorted in ascending order
func (it *Iterator[V]) IsSorted(less func(a, b V) bool) bool {
	prev := it.First()
	return it.All(func(v V) bool {
		defer func() {
			prev = v
		}()
		return !less(v, prev)
	})
}

// String returns a string representation of the Iterator
func (it *Iterator[V]) String() string {
	return fmt.Sprintf("<Iterator: %v>", it.Collect())
}

// Replace replaces all elements that satisfy the predicate with the replacement
func (it *Iterator[V]) Replace(predicate func(V) bool, replacement V) *Iterator[V] {
	return it.Map(func(v V) V {
		if predicate(v) {
			return replacement
		}
		return v
	})
}

// ReplaceAll replaces all elements with the replacement
func (it *Iterator[V]) ReplaceAll(replacement V) *Iterator[V] {
	return it.Replace(func(V) bool { return true }, replacement)
}

// Compact removes the nil elements from the Iterator
func (it *Iterator[V]) Compact() *Iterator[V] {
	return it.Filter(func(v V) bool { return !reflect.ValueOf(v).IsZero() })
}

// CompactWith removes the elements that are equal to the zero value of the type
func (it *Iterator[V]) CompactWith(zero V) *Iterator[V] {
	return it.Filter(func(v V) bool { return !reflect.DeepEqual(v, zero) })
}

// ToUpper converts all elements to uppercase if they are strings
// and if not, it leaves them unchanged in the Iterator
func (it *Iterator[V]) ToUpper() *Iterator[V] {
	return it.Map(func(v V) V {
		if str, ok := any(v).(string); ok {
			return any(strings.ToUpper(str)).(V)
		}
		return v
	})
}

// ToLower converts all elements to lowercase if they are strings
// and if not, it leaves them unchanged in the Iterator
func (it *Iterator[V]) ToLower() *Iterator[V] {
	return it.Map(func(v V) V {
		if str, ok := any(v).(string); ok {
			return any(strings.ToLower(str)).(V)
		}
		return v
	})
}

// TrimSpace trims the whitespace from all elements if they are strings
// and if not, it leaves them unchanged in the Iterator
func (it *Iterator[V]) TrimSpace() *Iterator[V] {
	return it.Map(func(v V) V {
		if str, ok := any(v).(string); ok {
			return any(strings.TrimSpace(str)).(V)
		}
		return v
	})
}
