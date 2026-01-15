// Package itertools provides a powerful and flexible generic iterator system for Go.
// It enables functional-style iteration, transformation, and aggregation of collections
// with support for lazy evaluation and composable operations.
//
// The library is inspired by iterator protocols found in Python and Rust, adapted
// for Go's type system and idioms using Go 1.23+ generics and iter.Seq.
package itertools

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"reflect"
	"sort"
	"strings"
)

// Iterator is a generic iterator that can iterate over any type of sequence.
// It supports both functional-style operations (via the seq field) and
// imperative-style iteration (via Next/Current methods).
//
// Iterator provides lazy evaluation - operations are only performed when
// elements are consumed. This allows efficient processing of large or
// infinite sequences.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	result := iter.
//	    Filter(func(x int) bool { return x%2 == 0 }).
//	    Map(func(x int) int { return x * x }).
//	    Collect()
//	// result is []int{4, 16}
type Iterator[V any] struct {
	seq  iter.Seq[V]
	curr *V
	done bool
	// Pull-based iterator for Next/Current
	pull func() (V, bool)
	stop func()
}

// ToIter creates an Iterator from a slice.
// The iterator will yield each element of the slice in order.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4})
//	for iter.Next() {
//	    fmt.Println(iter.Current())
//	}
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

// Next advances the iterator to the next element and returns true if successful.
// It returns false when the iterator is exhausted.
//
// Next must be called before accessing the current element via Current().
// After Next returns false, subsequent calls will continue to return false.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	for iter.Next() {
//	    fmt.Println(iter.Current())
//	}
func (it *Iterator[V]) Next() bool {
	if it.done {
		return false
	}

	// Lazy initialization: create pull iterator on first Next() call
	if it.pull == nil {
		it.pull, it.stop = iter.Pull(it.seq)
	}

	v, ok := it.pull()
	if !ok {
		it.done = true
		if it.stop != nil {
			it.stop()
		}
		return false
	}

	it.curr = &v
	return true
}

// Current returns the current element of the iterator.
// It panics if called before Next() or after the iterator is exhausted.
//
// Always call Next() before calling Current() to ensure a valid element exists.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	if iter.Next() {
//	    fmt.Println(iter.Current()) // Prints: 1
//	}
func (it *Iterator[V]) Current() V {
	if it.curr == nil {
		panic("iterator is not started")
	}

	if it.done {
		panic("iterator is done")
	}

	return *it.curr
}

// Close releases any resources held by the iterator, particularly cleaning up
// the goroutine spawned by iter.Pull() when using Next/Current iteration.
//
// IMPORTANT: Always call Close() when using Next/Current iteration and not
// fully exhausting the iterator. Failing to do so will leak a goroutine.
//
// Close is idempotent - it can be called multiple times safely.
// Close is automatically called when the iterator is fully exhausted.
//
// Best Practice - Use defer for guaranteed cleanup:
//
//	iter := itertools.Range(0, 1000000)
//	defer iter.Close()
//	for iter.Next() {
//	    if someCondition {
//	        break // Close() will be called via defer
//	    }
//	    process(iter.Current())
//	}
//
// Note: Close only affects Next/Current iteration. Functional-style operations
// (Collect, Each, Filter, Map, etc.) handle cleanup automatically.
func (it *Iterator[V]) Close() {
	if it.stop != nil {
		it.stop()
		it.stop = nil // Make Close idempotent
	}
	it.done = true
}

// Collect consumes the iterator and returns all elements as a slice.
// After calling Collect, the iterator is exhausted.
//
// Note: This loads all elements into memory. For large sequences,
// consider using streaming operations like Each or Filter.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4})
//	evens := iter.Filter(func(x int) bool { return x%2 == 0 }).Collect()
//	// evens is []int{2, 4}
func (it *Iterator[V]) Collect() []V {
	collect := make([]V, 0)
	it.seq(func(e V) bool {
		collect = append(collect, e)
		return true
	})
	return collect
}

// Each applies a function to each element of the iterator.
// This is useful for performing side effects like printing or logging.
//
// The iterator is consumed after calling Each.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	iter.Each(func(x int) {
//	    fmt.Println(x)
//	})
func (it *Iterator[V]) Each(f func(V)) {
	it.seq(func(v V) bool {
		f(v)
		return true
	})
}

// Reverse returns a new iterator that yields elements in reverse order.
//
// Note: This method collects all elements into memory to reverse them,
// so it's not suitable for infinite iterators.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4})
//	reversed := iter.Reverse().Collect()
//	// reversed is []int{4, 3, 2, 1}
func (it *Iterator[V]) Reverse() *Iterator[V] {
	xs := it.Collect()
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
	return ToIter(xs)
}

// Filter returns a new iterator that only yields elements satisfying the predicate.
// Elements for which predicate returns false are skipped.
//
// Filter is lazy - predicates are evaluated only as elements are consumed.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5, 6})
//	evens := iter.Filter(func(x int) bool { return x%2 == 0 }).Collect()
//	// evens is []int{2, 4, 6}
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

// Map transforms each element using the provided function and returns a new iterator.
//
// Map is lazy - the transformation is applied only as elements are consumed.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4})
//	squared := iter.Map(func(x int) int { return x * x }).Collect()
//	// squared is []int{1, 4, 9, 16}
func (it *Iterator[V]) Map(f func(V) V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			it.seq(func(v V) bool {
				return yield(f(v))
			})
		},
	}
}

// Chain concatenates this iterator with another, yielding all elements from
// the first iterator followed by all elements from the second.
//
// Chain properly handles early termination - if iteration stops early,
// the second iterator may not be consumed at all.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3})
//	iter2 := itertools.ToIter([]int{4, 5, 6})
//	chained := iter1.Chain(iter2).Collect()
//	// chained is []int{1, 2, 3, 4, 5, 6}
func (it *Iterator[V]) Chain(other *Iterator[V]) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
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
			other.seq(yield)
		},
	}
}

// Take returns a new iterator that yields at most the first n elements.
// If the iterator has fewer than n elements, all elements are yielded.
//
// If n is negative, it is treated as 0 (returns an empty iterator).
//
// Take is lazy and supports early termination.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	first3 := iter.Take(3).Collect()
//	// first3 is []int{1, 2, 3}
func (it *Iterator[V]) Take(n int) *Iterator[V] {
	if n < 0 {
		n = 0
	}
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

// Drop returns a new iterator that skips the first n elements and yields the rest.
// If the iterator has n or fewer elements, the resulting iterator is empty.
//
// If n is negative, it is treated as 0 (yields all elements).
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	rest := iter.Drop(2).Collect()
//	// rest is []int{3, 4, 5}
func (it *Iterator[V]) Drop(n int) *Iterator[V] {
	if n < 0 {
		n = 0
	}
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

// TakeWhile returns a new iterator that yields elements while the predicate is true.
// Once the predicate returns false for an element, iteration stops.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5, 1, 2})
//	result := iter.TakeWhile(func(x int) bool { return x < 4 }).Collect()
//	// result is []int{1, 2, 3}
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

// DropWhile returns a new iterator that skips elements while the predicate is true,
// then yields all remaining elements.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	result := iter.DropWhile(func(x int) bool { return x < 3 }).Collect()
//	// result is []int{3, 4, 5}
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

// First returns the first element of the iterator.
// It panics if the iterator is empty.
//
// For a safe alternative that doesn't panic, use FirstOr.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	first := iter.First() // Returns 1
func (it *Iterator[V]) First() V {
	it.Next()
	return it.Current()
}

// Last returns the last element of the iterator.
// It panics if the iterator is empty.
//
// For a safe alternative that doesn't panic, use LastOr.
//
// Note: This method collects all elements into memory.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	last := iter.Last() // Returns 3
func (it *Iterator[V]) Last() V {
	return it.Reverse().First()
}

// Nth returns the nth element (0-indexed) of the iterator.
// It panics if the iterator has fewer than n+1 elements.
//
// If n is negative, it is treated as 0 (returns the first element).
//
// For a safe alternative that doesn't panic, use NthOr.
//
// Example:
//
//	iter := itertools.ToIter([]int{10, 20, 30, 40})
//	third := iter.Nth(2) // Returns 30
func (it *Iterator[V]) Nth(n int) V {
	if n < 0 {
		n = 0
	}
	return it.Drop(n).First()
}

// FirstOr returns the first element of the iterator, or defaultValue if the iterator is empty.
// This is a safe alternative to First that doesn't panic.
//
// Example:
//
//	iter := itertools.ToIter([]int{})
//	first := iter.FirstOr(999) // Returns 999
//
//	iter2 := itertools.ToIter([]int{1, 2, 3})
//	first2 := iter2.FirstOr(999) // Returns 1
func (it *Iterator[V]) FirstOr(defaultValue V) V {
	var result V
	found := false
	it.seq(func(v V) bool {
		result = v
		found = true
		return false
	})
	if found {
		return result
	}
	return defaultValue
}

// LastOr returns the last element of the iterator, or defaultValue if the iterator is empty.
// This is a safe alternative to Last that doesn't panic.
//
// Note: This method must consume the entire iterator.
//
// Example:
//
//	iter := itertools.ToIter([]int{})
//	last := iter.LastOr(999) // Returns 999
//
//	iter2 := itertools.ToIter([]int{1, 2, 3})
//	last2 := iter2.LastOr(999) // Returns 3
func (it *Iterator[V]) LastOr(defaultValue V) V {
	var result V
	found := false
	it.seq(func(v V) bool {
		result = v
		found = true
		return true
	})
	if found {
		return result
	}
	return defaultValue
}

// NthOr returns the nth element (0-indexed) of the iterator, or defaultValue if there aren't enough elements.
// This is a safe alternative to Nth that doesn't panic.
//
// If n is negative, it is treated as 0 (returns the first element, or defaultValue if empty).
//
// Example:
//
//	iter := itertools.ToIter([]int{10, 20, 30})
//	third := iter.NthOr(2, 999)  // Returns 30
//	fifth := iter.NthOr(4, 999)  // Returns 999 (not enough elements)
func (it *Iterator[V]) NthOr(n int, defaultValue V) V {
	if n < 0 {
		n = 0
	}
	var result V
	found := false
	i := 0
	it.seq(func(v V) bool {
		if i == n {
			result = v
			found = true
			return false
		}
		i++
		return true
	})
	if found {
		return result
	}
	return defaultValue
}

// All returns true if all elements in the iterator satisfy the predicate.
// Returns true for an empty iterator.
//
// Example:
//
//	iter := itertools.ToIter([]int{2, 4, 6, 8})
//	allEven := iter.All(func(x int) bool { return x%2 == 0 }) // Returns true
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

// Any returns true if any element in the iterator satisfies the predicate.
// Returns false for an empty iterator.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 3, 5, 6, 7})
//	hasEven := iter.Any(func(x int) bool { return x%2 == 0 }) // Returns true
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

// Find returns the first element that satisfies the predicate along with true,
// or the zero value and false if no element satisfies the predicate.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 3, 5, 6, 7})
//	first, found := iter.Find(func(x int) bool { return x%2 == 0 })
//	// first is 6, found is true
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

// Sort returns a new iterator with elements sorted according to the less function.
//
// Note: This method collects all elements into memory to sort them.
//
// Example:
//
//	iter := itertools.ToIter([]int{3, 1, 4, 1, 5})
//	sorted := iter.Sort(func(a, b int) bool { return a < b }).Collect()
//	// sorted is []int{1, 1, 3, 4, 5}
func (it *Iterator[V]) Sort(less func(a, b V) bool) *Iterator[V] {
	xs := it.Collect()
	sort.Slice(xs, func(i, j int) bool {
		return less(xs[i], xs[j])
	})
	return ToIter(xs)
}

// Min returns the minimum element according to the less function, along with true.
// Returns the zero value and false if the iterator is empty.
//
// Example:
//
//	iter := itertools.ToIter([]int{3, 1, 4, 1, 5})
//	min, found := iter.Min(func(a, b int) bool { return a < b })
//	// min is 1, found is true
func (it *Iterator[V]) Min(less func(a, b V) bool) (V, bool) {
	var minValue V
	var found bool

	it.seq(func(v V) bool {
		if !found || less(v, minValue) {
			minValue = v
			found = true
		}
		return true
	})

	return minValue, found
}

// Max returns the maximum element according to the less function, along with true.
// Returns the zero value and false if the iterator is empty.
//
// Example:
//
//	iter := itertools.ToIter([]int{3, 1, 4, 1, 5})
//	max, found := iter.Max(func(a, b int) bool { return a < b })
//	// max is 5, found is true
func (it *Iterator[V]) Max(less func(a, b V) bool) (V, bool) {
	var maxValue V
	var found bool

	it.seq(func(v V) bool {
		if !found || less(maxValue, v) {
			maxValue = v
			found = true
		}
		return true
	})

	return maxValue, found
}

// Partition splits the iterator into two iterators: one with elements that satisfy
// the predicate, and one with elements that don't.
//
// Note: This method collects all elements into memory.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5, 6})
//	evens, odds := iter.Partition(func(x int) bool { return x%2 == 0 })
//	// evens.Collect() is []int{2, 4, 6}
//	// odds.Collect() is []int{1, 3, 5}
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

// Count returns the total number of elements in the iterator.
// The iterator is consumed after calling Count.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	count := iter.Filter(func(x int) bool { return x%2 == 0 }).Count()
//	// count is 2
func (it *Iterator[V]) Count() int {
	var count int
	it.seq(func(V) bool {
		count++
		return true
	})
	return count
}

// Unique returns a new iterator that yields only unique elements based on the key function.
// The keyFunc extracts a comparable key from each element.
//
// Each iteration creates a fresh set for tracking seen elements, making it safe
// to iterate multiple times.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 2, 3, 3, 3, 4})
//	unique := iter.Unique(func(x int) any { return x }).Collect()
//	// unique is []int{1, 2, 3, 4}
func (it *Iterator[V]) Unique(keyFunc func(V) any) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			seen := make(map[any]struct{}) // Fresh map for each iteration
			it.seq(func(v V) bool {
				key := keyFunc(v)
				if _, ok := seen[key]; ok {
					return true
				}
				seen[key] = struct{}{}
				return yield(v)
			})
		},
	}
}

// GroupBy groups elements by a key function into a map where keys are strings
// and values are slices of elements with that key.
//
// Example:
//
//	type Person struct { Name string; Age int }
//	iter := itertools.ToIter([]Person{
//	    {"Alice", 25}, {"Bob", 30}, {"Charlie", 25},
//	})
//	byAge := iter.GroupBy(func(p Person) string {
//	    return fmt.Sprintf("%d", p.Age)
//	})
//	// byAge["25"] is []Person{{"Alice", 25}, {"Charlie", 25}}
//	// byAge["30"] is []Person{{"Bob", 30}}
func (it *Iterator[V]) GroupBy(keyFunc func(V) string) map[string][]V {
	groups := make(map[string][]V)
	it.seq(func(v V) bool {
		key := keyFunc(v)
		groups[key] = append(groups[key], v)
		return true
	})
	return groups
}

// Cycle returns an infinite iterator that repeatedly cycles through the elements.
//
// Note: This collects all elements into memory. Be careful when using with
// infinite iterators or very large sequences.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	cycled := iter.Cycle().Take(7).Collect()
//	// cycled is []int{1, 2, 3, 1, 2, 3, 1}
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

// Repeat creates an iterator that yields the same value n times.
//
// Example:
//
//	iter := itertools.Repeat(42, 5)
//	result := iter.Collect()
//	// result is []int{42, 42, 42, 42, 42}
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

// Union returns an iterator that yields elements from both iterators without duplicates.
// The keyFunc extracts a comparable key to detect duplicates.
//
// Note: This method consumes both iterators and tracks seen keys in memory.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3})
//	iter2 := itertools.ToIter([]int{3, 4, 5})
//	union := iter1.Union(iter2, func(x int) any { return x }).Collect()
//	// union is []int{1, 2, 3, 4, 5}
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

// Intersection returns an iterator that yields elements present in both iterators.
// The keyFunc extracts a comparable key for matching elements.
//
// Note: This method consumes the other iterator to build a set of keys.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3, 4})
//	iter2 := itertools.ToIter([]int{3, 4, 5, 6})
//	intersection := iter1.Intersection(iter2, func(x int) any { return x }).Collect()
//	// intersection is []int{3, 4}
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

// Difference returns an iterator that yields elements present in this iterator
// but not in the other iterator.
// The keyFunc extracts a comparable key for matching elements.
//
// Note: This method consumes the other iterator to build a set of keys.
//
// Example:
//
//	iter1 := itertools.ToIter([]int{1, 2, 3, 4})
//	iter2 := itertools.ToIter([]int{3, 4, 5, 6})
//	diff := iter1.Difference(iter2, func(x int) any { return x }).Collect()
//	// diff is []int{1, 2}
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

// StepBy returns an iterator that yields every nth element (0-indexed).
//
// Panics if n <= 0.
//
// Example:
//
//	iter := itertools.ToIter([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
//	everyThird := iter.StepBy(3).Collect()
//	// everyThird is []int{0, 3, 6, 9}
func (it *Iterator[V]) StepBy(n int) *Iterator[V] {
	if n <= 0 {
		panic(fmt.Sprintf("StepBy: step must be positive, got %d", n))
	}
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

// Shuffle returns an iterator that yields elements in random order.
//
// Note: This collects all elements into memory to shuffle them.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	shuffled := iter.Shuffle().Collect()
//	// shuffled is []int in random order, e.g., []int{3, 1, 5, 2, 4}
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

// Index returns the 0-based index of the first element that satisfies the predicate.
// Returns -1 if no element satisfies the predicate.
//
// Example:
//
//	iter := itertools.ToIter([]int{10, 20, 30, 40, 50})
//	idx := iter.Index(func(x int) bool { return x == 30 })
//	// idx is 2
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

// LastIndex returns the 0-based index of the last element that satisfies the predicate.
// Returns -1 if no element satisfies the predicate.
//
// Example:
//
//	iter := itertools.ToIter([]int{10, 20, 30, 20, 50})
//	idx := iter.LastIndex(func(x int) bool { return x == 20 })
//	// idx is 3
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

// IsSorted returns true if the elements are sorted according to the less function.
// Returns true for empty iterators and single-element iterators.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	sorted := iter.IsSorted(func(a, b int) bool { return a < b })
//	// sorted is true
func (it *Iterator[V]) IsSorted(less func(a, b V) bool) bool {
	var prev V
	first := true
	result := true

	it.seq(func(v V) bool {
		if first {
			first = false
			prev = v
			return true
		}
		if less(v, prev) {
			result = false
			return false
		}
		prev = v
		return true
	})

	return result
}

// String returns a string representation of the iterator by collecting all elements.
//
// Note: This consumes the iterator and loads all elements into memory.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	str := iter.String()
//	// str is "<Iterator: [1 2 3]>"
func (it *Iterator[V]) String() string {
	return fmt.Sprintf("<Iterator: %v>", it.Collect())
}

// Replace returns an iterator that replaces elements satisfying the predicate
// with the replacement value.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	replaced := iter.Replace(func(x int) bool { return x%2 == 0 }, 0).Collect()
//	// replaced is []int{1, 0, 3, 0, 5}
func (it *Iterator[V]) Replace(predicate func(V) bool, replacement V) *Iterator[V] {
	return it.Map(func(v V) V {
		if predicate(v) {
			return replacement
		}
		return v
	})
}

// ReplaceAll returns an iterator that replaces all elements with the replacement value.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
//	replaced := iter.ReplaceAll(0).Collect()
//	// replaced is []int{0, 0, 0, 0, 0}
func (it *Iterator[V]) ReplaceAll(replacement V) *Iterator[V] {
	return it.Replace(func(V) bool { return true }, replacement)
}

// Compact removes zero-value elements from the iterator.
// Zero values are detected using reflection.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 0, 2, 0, 3})
//	compacted := iter.Compact().Collect()
//	// compacted is []int{1, 2, 3}
func (it *Iterator[V]) Compact() *Iterator[V] {
	return it.Filter(func(v V) bool { return !reflect.ValueOf(v).IsZero() })
}

// CompactWith removes elements that are equal to the specified zero value.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, -1, 2, -1, 3})
//	compacted := iter.CompactWith(-1).Collect()
//	// compacted is []int{1, 2, 3}
func (it *Iterator[V]) CompactWith(zero V) *Iterator[V] {
	return it.Filter(func(v V) bool { return !reflect.DeepEqual(v, zero) })
}

// ToUpper converts string elements to uppercase. Non-string elements are unchanged.
//
// Example:
//
//	iter := itertools.ToIter([]string{"hello", "world"})
//	upper := iter.ToUpper().Collect()
//	// upper is []string{"HELLO", "WORLD"}
func (it *Iterator[V]) ToUpper() *Iterator[V] {
	return it.Map(func(v V) V {
		if str, ok := any(v).(string); ok {
			return any(strings.ToUpper(str)).(V)
		}
		return v
	})
}

// ToLower converts string elements to lowercase. Non-string elements are unchanged.
//
// Example:
//
//	iter := itertools.ToIter([]string{"HELLO", "WORLD"})
//	lower := iter.ToLower().Collect()
//	// lower is []string{"hello", "world"}
func (it *Iterator[V]) ToLower() *Iterator[V] {
	return it.Map(func(v V) V {
		if str, ok := any(v).(string); ok {
			return any(strings.ToLower(str)).(V)
		}
		return v
	})
}

// TrimSpace trims whitespace from string elements. Non-string elements are unchanged.
//
// Example:
//
//	iter := itertools.ToIter([]string{"  hello  ", "  world  "})
//	trimmed := iter.TrimSpace().Collect()
//	// trimmed is []string{"hello", "world"}
func (it *Iterator[V]) TrimSpace() *Iterator[V] {
	return it.Map(func(v V) V {
		if str, ok := any(v).(string); ok {
			return any(strings.TrimSpace(str)).(V)
		}
		return v
	})
}

// AssertEq checks if the iterator elements are equal to the expected slice
// using the provided equality predicate.
// Returns true if all elements match, false otherwise.
//
// Example:
//
//	iter := itertools.ToIter([]int{1, 2, 3})
//	equal := iter.AssertEq([]int{1, 2, 3}, func(a, b int) bool { return a == b })
//	// equal is true
func (it *Iterator[V]) AssertEq(expected []V, predicate func(V, V) bool) bool {
	actual := it.Collect()
	if len(actual) != len(expected) {
		return false
	}
	for i, v := range actual {
		if !predicate(v, expected[i]) {
			return false
		}
	}
	return true
}
