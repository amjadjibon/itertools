package itertools_test

import (
	"testing"

	"github.com/amjadjibon/itertools"
)

// BenchmarkIterator_Next_TakeFirst tests taking only the first element
// This demonstrates lazy evaluation - we shouldn't process all elements
func BenchmarkIterator_Next_TakeFirst(b *testing.B) {
	slice := make([]int, 1000000) // 1 million elements
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.ToIter(slice)
		if iter.Next() {
			_ = iter.Current()
		}
	}
}

// BenchmarkIterator_Next_TakeAll tests iterating through all elements
func BenchmarkIterator_Next_TakeAll(b *testing.B) {
	slice := make([]int, 10000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.ToIter(slice)
		for iter.Next() {
			_ = iter.Current()
		}
	}
}

// BenchmarkIterator_Collect tests the functional style
func BenchmarkIterator_Collect(b *testing.B) {
	slice := make([]int, 10000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.ToIter(slice)
		_ = iter.Collect()
	}
}

// BenchmarkIterator_Filter_TakeFirst tests lazy evaluation with filters
func BenchmarkIterator_Filter_TakeFirst(b *testing.B) {
	slice := make([]int, 1000000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.ToIter(slice).Filter(func(v int) bool {
			return v%2 == 0
		})
		if iter.Next() {
			_ = iter.Current()
		}
	}
}

// BenchmarkIterator_Chain tests chaining operations
func BenchmarkIterator_Chain(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.ToIter(slice).
			Filter(func(v int) bool { return v%2 == 0 }).
			Map(func(v int) int { return v * 2 }).
			Take(10)
		_ = iter.Collect()
	}
}
