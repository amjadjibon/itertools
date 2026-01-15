package itertools_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/amjadjibon/itertools"
	"github.com/stretchr/testify/assert"
)

func TestNewIterator(t *testing.T) {
	iter1 := itertools.NewIterator(1, 2, 3, 4, 5)
	collected1 := iter1.Collect()

	assert.Equal(t, 5, len(collected1))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, collected1)

	iter2 := itertools.NewIterator("a", "b", "c", "d", "e")
	collected2 := iter2.Collect()

	assert.Equal(t, 5, len(collected2))
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, collected2)
}

func TestZip(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5, 6}
	slice2 := []string{"a", "b", "c", "d", "e"}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)
	zip := itertools.Zip(iter1, iter2).Collect()

	if len(zip) != 5 {
		t.Errorf("expected 5 elements, got %d", len(zip))
	}

	for i, v := range zip {
		if slice1[i] != v.First {
			t.Errorf("expected %d, got %d", v.First, slice1[i])
		}
		if slice2[i] != v.Second {
			t.Errorf("expected %s, got %s", v.Second, slice2[i])
		}
	}
}

func TestZip2(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []string{"a", "b", "c"}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)
	zip := itertools.Zip2(iter1, iter2, struct {
		First  int
		Second string
	}{0, ""}).Collect()

	if len(zip) != 5 {
		t.Errorf("expected 5 elements, got %d", len(zip))
	}

	for i, v := range zip {
		if i < len(slice1) {
			if slice1[i] != v.First {
				t.Errorf("expected %d, got %d", v.First, slice1[i])
			}
		} else {
			if v.First != 0 {
				t.Errorf("expected 0, got %d", v.First)
			}
		}

		if i < len(slice2) {
			if slice2[i] != v.Second {
				t.Errorf("expected %s, got %s", v.Second, slice2[i])
			}
		} else {
			if v.Second != "" {
				t.Errorf("expected \"\", got %s", v.Second)
			}
		}
	}
}

func TestSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	sum := itertools.Sum(iter, func(v int) int { return v }, 0)
	assert.Equal(t, 15, sum)
}

func TestSumFloat(t *testing.T) {
	slice := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	iter := itertools.ToIter(slice)

	sum := itertools.Sum(iter, func(v float64) float64 { return v }, 0)
	assert.Equal(t, 16.5, sum)
}

func TestSumComplex(t *testing.T) {
	type Complex struct {
		A int
		B int
	}

	slice := []Complex{{1, 2}, {3, 4}, {5, 6}}
	iter := itertools.ToIter(slice)

	sum := itertools.Sum(iter, func(v Complex) int { return v.A + v.B }, 0)
	assert.Equal(t, 21, sum)
}

func TestFoldSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	sum := itertools.Fold(iter, func(acc, v int) int { return acc + v }, 0)
	assert.Equal(t, 15, sum)
}

func TestFoldProduct(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	product := itertools.Fold(iter, func(acc, v int) int { return acc * v }, 1)
	assert.Equal(t, 120, product)
}

func TestFoldConcat(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e"}
	iter := itertools.ToIter(slice)

	concat := itertools.Fold(iter, func(acc, v string) string { return acc + v }, "")
	assert.Equal(t, "abcde", concat)
}

func TestProduct(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	product := itertools.Product(iter, func(v int) int { return v }, 1)
	assert.Equal(t, 120, product)
}

func TestProductFloat(t *testing.T) {
	slice := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	iter := itertools.ToIter(slice)

	product := itertools.Product(iter, func(v float64) float64 { return v }, 1)
	assert.Equal(t, fmt.Sprintf("%.2f", 1.1*2.2*3.3*4.4*5.5), fmt.Sprintf("%.2f", product))
}

func TestProductComplex(t *testing.T) {
	type Complex struct {
		A int
		B int
	}

	slice := []Complex{{1, 2}, {3, 4}, {5, 6}}
	iter := itertools.ToIter(slice)

	product := itertools.Product(iter, func(v Complex) int { return v.A * v.B }, 1)
	assert.Equal(t, 720, product)
}

func TestChunkSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	iter := itertools.ToIter(slice)

	chunks := itertools.ChunkSlice(iter, 3)
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}}

	for i := 0; chunks.Next(); i++ {
		chunk := itertools.ToIter(chunks.Current())
		for j := 0; chunk.Next(); j++ {
			assert.Equal(t, expected[i][j], chunk.Current())
		}
	}
}

func TestLazyChunks(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	iter := itertools.ToIter(slice)

	chunks := itertools.Chunks(iter, 3)
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}}

	for i := 0; chunks.Next(); i++ {
		chunk := chunks.Current()
		for j := 0; chunk.Next(); j++ {
			assert.Equal(t, expected[i][j], chunk.Current())
		}
	}
}

func TestFlatten(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	iter := itertools.ToIter(slice)

	chunkList := itertools.ChunkList(iter, 3)

	flatten := itertools.Flatten(chunkList...)
	assert.Equal(t, slice, flatten.Collect())
}

func TestCartesianProduct(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []string{"a", "b", "c"}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	cartesianProduct := itertools.CartesianProduct(iter1, iter2).Collect()
	expected := []struct {
		X int
		Y string
	}{
		{1, "a"},
		{1, "b"},
		{1, "c"},
		{2, "a"},
		{2, "b"},
		{2, "c"},
		{3, "a"},
		{3, "b"},
		{3, "c"},
	}

	for i, v := range cartesianProduct {
		assert.Equal(t, expected[i].X, v.X)
		assert.Equal(t, expected[i].Y, v.Y)
	}
}

func TestCartesianProductEmpty(t *testing.T) {
	slice1 := []int{}
	slice2 := []string{}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	cartesianProduct := itertools.CartesianProduct(iter1, iter2).Collect()
	assert.Empty(t, cartesianProduct)
}

// TestZip_NoGoroutineLeak verifies that Zip properly cleans up goroutines
// when the iterator is stopped early
func TestZip_NoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	// Create large iterators but only consume 5 elements
	iter1 := itertools.Range(0, 1000000)
	iter2 := itertools.Range(0, 1000000)

	zipped := itertools.Zip(iter1, iter2).Take(5).Collect()

	assert.Equal(t, 5, len(zipped))

	// Give time for goroutines to clean up
	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()

	// Should not have leaked goroutines
	assert.LessOrEqual(t, after, before+1, "Goroutine leak detected")
}

// TestZip2_NoGoroutineLeak verifies that Zip2 properly cleans up goroutines
func TestZip2_NoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	iter1 := itertools.Range(0, 1000000)
	iter2 := itertools.Range(0, 1000000)

	fill := struct {
		First  int
		Second int
	}{-1, -1}

	zipped := itertools.Zip2(iter1, iter2, fill).Take(5).Collect()

	assert.Equal(t, 5, len(zipped))

	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()

	assert.LessOrEqual(t, after, before+1, "Goroutine leak detected")
}

// TestZip2_FillValues verifies that Zip2 actually uses fill values
func TestZip2_FillValues(t *testing.T) {
	iter1 := itertools.ToIter([]int{1, 2, 3, 4, 5})
	iter2 := itertools.ToIter([]string{"a", "b"})

	fill := struct {
		First  int
		Second string
	}{-1, "FILL"}

	result := itertools.Zip2(iter1, iter2, fill).Collect()

	assert.Equal(t, 5, len(result))
	assert.Equal(t, 1, result[0].First)
	assert.Equal(t, "a", result[0].Second)
	assert.Equal(t, 2, result[1].First)
	assert.Equal(t, "b", result[1].Second)
	// After iter2 ends, should use fill value
	assert.Equal(t, 3, result[2].First)
	assert.Equal(t, "FILL", result[2].Second)
	assert.Equal(t, 4, result[3].First)
	assert.Equal(t, "FILL", result[3].Second)
	assert.Equal(t, 5, result[4].First)
	assert.Equal(t, "FILL", result[4].Second)
}

// TestZip2_FillBothSides verifies fill works when iter1 is shorter
func TestZip2_FillBothSides(t *testing.T) {
	iter1 := itertools.ToIter([]int{1, 2})
	iter2 := itertools.ToIter([]string{"a", "b", "c", "d"})

	fill := struct {
		First  int
		Second string
	}{-99, "EMPTY"}

	result := itertools.Zip2(iter1, iter2, fill).Collect()

	assert.Equal(t, 4, len(result))
	assert.Equal(t, 1, result[0].First)
	assert.Equal(t, "a", result[0].Second)
	assert.Equal(t, 2, result[1].First)
	assert.Equal(t, "b", result[1].Second)
	// After iter1 ends, should use fill value
	assert.Equal(t, -99, result[2].First)
	assert.Equal(t, "c", result[2].Second)
	assert.Equal(t, -99, result[3].First)
	assert.Equal(t, "d", result[3].Second)
}

// TestFlatten_EarlyTermination verifies that Flatten stops properly
func TestFlatten_EarlyTermination(t *testing.T) {
	iter1 := itertools.Range(0, 100)
	iter2 := itertools.Range(100, 200)
	iter3 := itertools.Range(200, 300)

	// Take only 5 elements - should not process iter2 or iter3
	result := itertools.Flatten(iter1, iter2, iter3).Take(5).Collect()

	assert.Equal(t, 5, len(result))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
}

// TestChunkList_Functional verifies ChunkList works correctly
func TestChunkList_Functional(t *testing.T) {
	iter := itertools.Range(0, 10)
	chunks := itertools.ChunkList(iter, 3)

	assert.Equal(t, 4, len(chunks))
	assert.Equal(t, []int{0, 1, 2}, chunks[0].Collect())
	assert.Equal(t, []int{3, 4, 5}, chunks[1].Collect())
	assert.Equal(t, []int{6, 7, 8}, chunks[2].Collect())
	assert.Equal(t, []int{9}, chunks[3].Collect())
}

// =============================================================================
// PANIC AND EDGE CASE TESTS
// =============================================================================

// TestChunkSlice_PanicOnNegativeSize verifies ChunkSlice panics with negative size
func TestChunkSlice_PanicOnNegativeSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	assert.Panics(t, func() {
		itertools.ChunkSlice(iter, -5).Collect()
	}, "ChunkSlice should panic with negative size")
}

// TestChunks_PanicOnNegativeSize verifies Chunks panics with negative size
func TestChunks_PanicOnNegativeSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	assert.Panics(t, func() {
		itertools.Chunks(iter, -3).Collect()
	}, "Chunks should panic with negative size")
}

// TestChunkList_PanicOnNegativeSize verifies ChunkList panics with negative size
func TestChunkList_PanicOnNegativeSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	assert.Panics(t, func() {
		itertools.ChunkList(iter, -2)
	}, "ChunkList should panic with negative size")
}

// TestChunkSlice_ZeroSize tests that ChunkSlice panics with size=0
func TestChunkSlice_ZeroSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Now correctly panics with zero size
	assert.Panics(t, func() {
		itertools.ChunkSlice(iter, 0).Collect()
	}, "ChunkSlice should panic with zero size")
}

// TestChunks_ZeroSize tests that Chunks panics with size=0
func TestChunks_ZeroSize(t *testing.T) {
	iter := itertools.Range(0, 10)

	// Now correctly panics with zero size
	assert.Panics(t, func() {
		itertools.Chunks(iter, 0).Collect()
	}, "Chunks should panic with zero size")
}

// TestFold_BasicOperation tests Fold with basic addition
func TestFold_BasicOperation(t *testing.T) {
	iter := itertools.Range(1, 6) // 1, 2, 3, 4, 5

	sum := itertools.Fold(iter, func(acc, v int) int {
		return acc + v
	}, 0)

	assert.Equal(t, 15, sum)
}

// TestFold_StringConcatenation tests Fold with string concatenation
func TestFold_StringConcatenation(t *testing.T) {
	iter := itertools.ToIter([]string{"a", "b", "c"})

	result := itertools.Fold(iter, func(acc, v string) string {
		return acc + v
	}, "")

	assert.Equal(t, "abc", result)
}

// TestCartesianProduct_EmptyIterators tests CartesianProduct with empty inputs
func TestCartesianProduct_EmptyIterators(t *testing.T) {
	iter1 := itertools.ToIter([]int{})
	iter2 := itertools.ToIter([]string{"a", "b"})

	result := itertools.CartesianProduct(iter1, iter2).Collect()

	assert.Empty(t, result, "CartesianProduct with empty iterator should return empty")
}

// TestCartesianProduct_SingleElements tests CartesianProduct with single elements
func TestCartesianProduct_SingleElements(t *testing.T) {
	iter1 := itertools.ToIter([]int{1})
	iter2 := itertools.ToIter([]string{"x"})

	result := itertools.CartesianProduct(iter1, iter2).Collect()

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 1, result[0].X)
	assert.Equal(t, "x", result[0].Y)
}
