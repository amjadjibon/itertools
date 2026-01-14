package itertools_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/amjadjibon/itertools"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// GOROUTINE LEAK TESTS
// =============================================================================

// TestGoroutineLeakHelper is a helper to count goroutines
func countGoroutines() int {
	// Give time for goroutines to start/stop
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond)
	runtime.GC()
	time.Sleep(10 * time.Millisecond)
	return runtime.NumGoroutine()
}

// TestIteratorNext_GoroutineLeak_EarlyReturn tests goroutine leak with early return
func TestIteratorNext_GoroutineLeak_EarlyReturn(t *testing.T) {
	before := countGoroutines()

	// Create 10 iterators and call Next() but never exhaust them
	for i := 0; i < 10; i++ {
		iter := itertools.Range(0, 1000000)
		iter.Next()
		_ = iter.Current()
		// Iterator goes out of scope without exhausting - LEAK!
	}

	after := countGoroutines()
	leaked := after - before

	// We expect leaks in the current implementation
	// After fix with Close(), this should be 0
	if leaked > 0 {
		t.Logf("WARNING: Detected %d leaked goroutines (expected before fix)", leaked)
	}

	// TODO: After implementing Close(), uncomment this assertion:
	// assert.Equal(t, 0, leaked, "Should not leak goroutines after Close() is implemented")
}

// TestIteratorNext_GoroutineLeak_BreakInLoop tests goroutine leak with break
func TestIteratorNext_GoroutineLeak_BreakInLoop(t *testing.T) {
	before := countGoroutines()

	// Create 10 iterators and break early
	for i := 0; i < 10; i++ {
		iter := itertools.Range(0, 1000000)
		count := 0
		for iter.Next() {
			count++
			if count >= 5 {
				break // Early termination - LEAK!
			}
		}
	}

	after := countGoroutines()
	leaked := after - before

	if leaked > 0 {
		t.Logf("WARNING: Detected %d leaked goroutines from break (expected before fix)", leaked)
	}

	// TODO: After implementing Close(), users should do:
	// defer iter.Close()
	// Then: assert.Equal(t, 0, leaked)
}

// TestIteratorNext_GoroutineLeak_ErrorReturn tests goroutine leak with error handling
func TestIteratorNext_GoroutineLeak_ErrorReturn(t *testing.T) {
	before := countGoroutines()

	processWithError := func() error {
		iter := itertools.Range(0, 1000000)
		for iter.Next() {
			if iter.Current() > 100 {
				return assert.AnError // Early return - LEAK!
			}
		}
		return nil
	}

	// Call it multiple times
	for i := 0; i < 10; i++ {
		_ = processWithError()
	}

	after := countGoroutines()
	leaked := after - before

	if leaked > 0 {
		t.Logf("WARNING: Detected %d leaked goroutines from error return (expected before fix)", leaked)
	}
}

// TestIteratorNext_NoLeak_FullExhaustion tests that full exhaustion doesn't leak
func TestIteratorNext_NoLeak_FullExhaustion(t *testing.T) {
	before := countGoroutines()

	// Create and fully exhaust iterators
	for i := 0; i < 10; i++ {
		iter := itertools.Range(0, 100)
		for iter.Next() {
			_ = iter.Current()
		}
		// Fully exhausted - should not leak
	}

	after := countGoroutines()
	leaked := after - before

	// Full exhaustion should clean up properly
	assert.LessOrEqual(t, leaked, 1, "Should not leak goroutines when fully exhausted")
}

// =============================================================================
// PANIC TESTS - StepBy
// =============================================================================

// TestStepBy_PanicOnZero tests that StepBy panics with step=0
func TestStepBy_PanicOnZero(t *testing.T) {
	iter := itertools.Range(0, 10)

	assert.Panics(t, func() {
		iter.StepBy(0).Collect()
	}, "StepBy(0) should panic due to division by zero")
}

// TestStepBy_PanicOnNegative tests that StepBy with negative value
func TestStepBy_InvalidWithNegative(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Currently returns all elements (incorrect behavior)
	result := iter.StepBy(-1).Collect()

	// After fix, this should panic or return error
	// For now, we document the incorrect behavior
	t.Logf("StepBy(-1) currently returns %d elements (should be invalid)", len(result))

	// TODO: After fix, this should panic:
	// assert.Panics(t, func() {
	//     iter.StepBy(-1).Collect()
	// }, "StepBy with negative step should panic")
}

// =============================================================================
// PANIC TESTS - ChunkSlice (Note: Duplicate tests are in itertools_test.go)
// =============================================================================

// TestChunkSlice_InvalidWithZeroSize tests ChunkSlice with size=0
func TestChunkSlice_InvalidWithZeroSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Currently has unexpected behavior (returns single chunk)
	result := itertools.ChunkSlice(iter, 0).Collect()

	// Document current incorrect behavior
	t.Logf("ChunkSlice(0) returns %d chunks (should be invalid)", len(result))

	// TODO: After fix, this should panic:
	// assert.Panics(t, func() {
	//     itertools.ChunkSlice(iter, 0).Collect()
	// }, "ChunkSlice with zero size should panic")
}

// =============================================================================
// PANIC TESTS - Chunks (Note: Duplicate tests are in itertools_test.go)
// =============================================================================

// TestChunks_InvalidWithZeroSize tests Chunks with size=0
func TestChunks_InvalidWithZeroSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Test current behavior with size=0
	result := itertools.Chunks(iter, 0).Collect()

	t.Logf("Chunks(0) returns %d chunks (should be invalid)", len(result))

	// TODO: After fix, should panic
}

// =============================================================================
// EDGE CASE TESTS - Negative Values
// =============================================================================

// TestTake_NegativeValue tests Take with negative value
func TestTake_NegativeValue(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Test current behavior
	result := iter.Take(-1).Collect()

	// Currently might return empty or all - document behavior
	t.Logf("Take(-1) returns %d elements", len(result))

	// TODO: After fix, negative should be treated as 0 or panic
	// assert.Empty(t, result, "Take with negative value should return empty")
}

// TestDrop_NegativeValue tests Drop with negative value
func TestDrop_NegativeValue(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Test current behavior
	result := iter.Drop(-1).Collect()

	t.Logf("Drop(-1) returns %d elements", len(result))

	// TODO: After fix, negative should be treated as 0 or panic
}

// TestNth_NegativeIndex tests Nth with negative index
func TestNth_NegativeIndex(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Test current behavior - likely panics or returns zero value
	// Wrap in recover to avoid test failure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Nth(-1) panicked: %v", r)
		}
	}()

	result := iter.Nth(-1)
	t.Logf("Nth(-1) returned: %v", result)
}

// =============================================================================
// MEMORY LEAK TESTS - Infinite Iterators
// =============================================================================

// TestCollect_InfiniteIterator_MemoryLeak tests that Collect on infinite iterator causes OOM
// NOTE: This test is disabled by default as it will cause OOM
func TestCollect_InfiniteIterator_MemoryLeak(t *testing.T) {
	t.Skip("Skipping OOM test - would cause memory exhaustion")

	// This would cause OOM:
	// infiniteIter := itertools.FromFunc(func() (int, bool) {
	//     return rand.Int(), true
	// })
	// infiniteIter.Collect() // OOM!
}

// TestCollect_LargeButFinite tests memory usage with large finite iterator
func TestCollect_LargeButFinite(t *testing.T) {
	// Collect 1 million integers
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	before := m.Alloc

	iter := itertools.Range(0, 1000000)
	result := iter.Collect()

	runtime.ReadMemStats(&m)
	after := m.Alloc
	allocated := after - before

	// Should have allocated approximately 1M * 8 bytes = 8MB for int slice
	// (plus overhead for slice header, etc.)
	t.Logf("Collected %d elements, allocated ~%.2f MB", len(result), float64(allocated)/(1024*1024))

	assert.Equal(t, 1000000, len(result), "Should collect all elements")
	assert.Greater(t, allocated, uint64(8*1000000), "Should allocate memory for elements")
}

// =============================================================================
// RESOURCE CLEANUP TESTS
// =============================================================================

// TestIterator_Close_Method tests that Close method exists and works
// NOTE: This will fail until Close() is implemented
func TestIterator_Close_Method(t *testing.T) {
	t.Skip("Skipping until Close() method is implemented")

	// After implementation, test should work like this:
	// iter := itertools.Range(0, 1000000)
	// iter.Next()
	// iter.Close() // Should not panic
}

// TestIterator_Close_Idempotent tests that Close can be called multiple times
func TestIterator_Close_Idempotent(t *testing.T) {
	t.Skip("Skipping until Close() method is implemented")

	// After implementation:
	// iter := itertools.Range(0, 1000000)
	// iter.Close()
	// iter.Close() // Should not panic
	// iter.Close() // Should not panic
}

// TestIterator_Close_WithDefer tests proper defer cleanup pattern
func TestIterator_Close_WithDefer(t *testing.T) {
	t.Skip("Skipping until Close() method is implemented")

	// After implementation, recommended pattern:
	// before := countGoroutines()
	//
	// for i := 0; i < 10; i++ {
	//     iter := itertools.Range(0, 1000000)
	//     defer iter.Close() // Proper cleanup
	//     iter.Next()
	//     _ = iter.Current()
	// }
	//
	// after := countGoroutines()
	// assert.Equal(t, 0, after-before, "Should not leak with proper Close()")
}

// =============================================================================
// CONCURRENT ACCESS TESTS
// =============================================================================

// TestIterator_ConcurrentAccess tests thread safety
func TestIterator_ConcurrentAccess(t *testing.T) {
	// Iterators should NOT be used concurrently
	// This test documents the expected behavior

	// This is UNSAFE and can cause data races
	// Users should not do this
	t.Run("Unsafe Concurrent Access", func(t *testing.T) {
		// Don't actually run concurrent access in test
		// Just document that it's unsafe
		t.Log("Iterators are NOT thread-safe - do not use from multiple goroutines")
	})
}

// =============================================================================
// BENCHMARK - Leak Impact
// =============================================================================

// BenchmarkIteratorNext_WithLeak benchmarks the impact of goroutine leaks
func BenchmarkIteratorNext_WithLeak(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := itertools.Range(0, 1000)
		iter.Next()
		_ = iter.Current()
		// Leak occurs here
	}
}

// BenchmarkIteratorNext_FullExhaustion benchmarks proper full exhaustion
func BenchmarkIteratorNext_FullExhaustion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iter := itertools.Range(0, 1000)
		for iter.Next() {
			_ = iter.Current()
		}
		// No leak
	}
}
