package itertools_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/amjadjibon/itertools"
	"github.com/stretchr/testify/assert"
)

func TestFromChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	iter := itertools.FromChannel(ch)
	result := iter.Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFromChannel_WithFilter(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	iter := itertools.FromChannel(ch)
	result := iter.Filter(func(x int) bool {
		return x%2 == 0
	}).Collect()

	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestFromChannel_EarlyTermination(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()

	iter := itertools.FromChannel(ch)
	result := iter.Take(5).Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFromChannelWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int)
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
			time.Sleep(1 * time.Millisecond)
		}
		close(ch)
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	iter := itertools.FromChannelWithContext(ctx, ch)
	result := iter.Collect()

	// Should have collected some elements before context was cancelled
	assert.Greater(t, len(result), 0)
	assert.Less(t, len(result), 100)
}

func TestFromReader(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5"
	reader := strings.NewReader(input)

	iter := itertools.FromReader(reader)
	result := iter.Collect()

	assert.Equal(t, []string{"line1", "line2", "line3", "line4", "line5"}, result)
}

func TestFromReader_WithFilter(t *testing.T) {
	input := "apple\nbanana\napricot\nblueberry\navocado"
	reader := strings.NewReader(input)

	iter := itertools.FromReader(reader)
	result := iter.Filter(func(line string) bool {
		return strings.HasPrefix(line, "a")
	}).Collect()

	assert.Equal(t, []string{"apple", "apricot", "avocado"}, result)
}

func TestFromReader_EarlyTermination(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"
	reader := strings.NewReader(input)

	iter := itertools.FromReader(reader)
	result := iter.Take(3).Collect()

	assert.Equal(t, []string{"line1", "line2", "line3"}, result)
}

func TestFromReaderWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	input := strings.Repeat("line\n", 1000000) // Much larger input
	reader := strings.NewReader(input)

	// Cancel almost immediately to ensure we don't read everything
	cancel()

	iter := itertools.FromReaderWithContext(ctx, reader)
	result := iter.Collect()

	// Should have stopped early due to cancelled context
	assert.Less(t, len(result), 1000000)
}

func TestFromFunc(t *testing.T) {
	counter := 0
	iter := itertools.FromFunc(func() (int, bool) {
		if counter >= 5 {
			return 0, false
		}
		counter++
		return counter, true
	})

	result := iter.Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFromFunc_Fibonacci(t *testing.T) {
	a, b := 0, 1
	iter := itertools.FromFunc(func() (int, bool) {
		result := a
		a, b = b, a+b
		return result, true
	})

	result := iter.Take(10).Collect()
	assert.Equal(t, []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}, result)
}

func TestFromFunc_EarlyTermination(t *testing.T) {
	counter := 0
	iter := itertools.FromFunc(func() (int, bool) {
		counter++
		return counter, true // infinite
	})

	result := iter.TakeWhile(func(x int) bool {
		return x <= 7
	}).Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, result)
}

func TestFromFuncWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	counter := 0
	iter := itertools.FromFuncWithContext(ctx, func() (int, bool) {
		counter++
		time.Sleep(1 * time.Millisecond)
		return counter, true
	})

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	result := iter.Collect()

	assert.Greater(t, len(result), 0)
	assert.Less(t, len(result), 100)
}

func TestRange(t *testing.T) {
	iter := itertools.Range(0, 5)
	result := iter.Collect()

	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
}

func TestRange_Negative(t *testing.T) {
	iter := itertools.Range(-3, 3)
	result := iter.Collect()

	assert.Equal(t, []int{-3, -2, -1, 0, 1, 2}, result)
}

func TestRange_WithMap(t *testing.T) {
	iter := itertools.Range(1, 6)
	result := iter.Map(func(x int) int {
		return x * x
	}).Collect()

	assert.Equal(t, []int{1, 4, 9, 16, 25}, result)
}

func TestRangeStep(t *testing.T) {
	iter := itertools.RangeStep(0, 10, 2)
	result := iter.Collect()

	assert.Equal(t, []int{0, 2, 4, 6, 8}, result)
}

func TestRangeStep_Negative(t *testing.T) {
	iter := itertools.RangeStep(10, 0, -2)
	result := iter.Collect()

	assert.Equal(t, []int{10, 8, 6, 4, 2}, result)
}

func TestRangeStep_LargeStep(t *testing.T) {
	iter := itertools.RangeStep(0, 100, 25)
	result := iter.Collect()

	assert.Equal(t, []int{0, 25, 50, 75}, result)
}

func TestGenerate(t *testing.T) {
	counter := 0
	iter := itertools.Generate(func() int {
		counter++
		return counter * 2
	})

	result := iter.Take(5).Collect()
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestGenerate_Constant(t *testing.T) {
	iter := itertools.Generate(func() string {
		return "hello"
	})

	result := iter.Take(3).Collect()
	assert.Equal(t, []string{"hello", "hello", "hello"}, result)
}

func TestGenerateWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	counter := 0
	iter := itertools.GenerateWithContext(ctx, func() int {
		counter++
		time.Sleep(5 * time.Millisecond)
		return counter
	})

	result := iter.Collect()

	// Should have generated some elements before timeout
	assert.Greater(t, len(result), 0)
	assert.Less(t, len(result), 100)
}

func TestStreamChaining(t *testing.T) {
	// Test complex chaining with stream sources
	ch := make(chan int)
	go func() {
		for i := 1; i <= 20; i++ {
			ch <- i
		}
		close(ch)
	}()

	result := itertools.FromChannel(ch).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Collect()

	assert.Equal(t, []int{4, 16, 36}, result)
}

func TestStreamWithNext(t *testing.T) {
	// Test that stream iterators work with Next/Current pattern
	iter := itertools.Range(1, 6)

	var result []int
	for iter.Next() {
		result = append(result, iter.Current())
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}
