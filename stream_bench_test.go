package itertools_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/amjadjibon/itertools"
)

func BenchmarkFromChannel(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := make(chan int, 100)
		go func() {
			for j := 0; j < 100; j++ {
				ch <- j
			}
			close(ch)
		}()

		iter := itertools.FromChannel(ch)
		_ = iter.Collect()
	}
}

func BenchmarkFromChannel_EarlyTermination(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := make(chan int, 1000000)
		go func() {
			for j := 0; j < 1000000; j++ {
				ch <- j
			}
			close(ch)
		}()

		iter := itertools.FromChannel(ch)
		_ = iter.Take(10).Collect()
	}
}

func BenchmarkFromReader(b *testing.B) {
	input := strings.Repeat("line\n", 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		iter := itertools.FromReader(reader)
		_ = iter.Collect()
	}
}

func BenchmarkFromReader_WithFilter(b *testing.B) {
	input := strings.Repeat("apple\nbanana\napricot\nblueberry\n", 250)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(input)
		iter := itertools.FromReader(reader)
		_ = iter.Filter(func(line string) bool {
			return strings.HasPrefix(line, "a")
		}).Collect()
	}
}

func BenchmarkFromFunc_Fibonacci(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a, b := 0, 1
		iter := itertools.FromFunc(func() (int, bool) {
			result := a
			a, b = b, a+b
			return result, true
		})
		_ = iter.Take(20).Collect()
	}
}

func BenchmarkRange(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.Range(0, 10000)
		_ = iter.Collect()
	}
}

func BenchmarkRange_WithMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.Range(0, 1000)
		_ = iter.Map(func(x int) int {
			return x * x
		}).Collect()
	}
}

func BenchmarkRangeStep(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := itertools.RangeStep(0, 10000, 10)
		_ = iter.Collect()
	}
}

func BenchmarkGenerate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter := 0
		iter := itertools.Generate(func() int {
			counter++
			return counter
		})
		_ = iter.Take(1000).Collect()
	}
}

func BenchmarkStreamChaining(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := make(chan int, 1000)
		go func() {
			for j := 0; j < 1000; j++ {
				ch <- j
			}
			close(ch)
		}()

		_ = itertools.FromChannel(ch).
			Filter(func(x int) bool { return x%2 == 0 }).
			Map(func(x int) int { return x * x }).
			Take(100).
			Collect()
	}
}

func BenchmarkFromChannelWithContext(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		ch := make(chan int, 100)
		go func() {
			for j := 0; j < 100; j++ {
				ch <- j
			}
			close(ch)
		}()

		iter := itertools.FromChannelWithContext(ctx, ch)
		_ = iter.Collect()
		cancel()
	}
}
