package itertools

import (
	"bufio"
	"context"
	"io"
)

// FromChannel creates a lazy Iterator from a channel.
// The iterator will consume elements from the channel until it's closed.
// This is useful for concurrent producer-consumer patterns.
//
// Example:
//
//	ch := make(chan int)
//	go func() {
//	    for i := 0; i < 10; i++ {
//	        ch <- i
//	    }
//	    close(ch)
//	}()
//	iter := itertools.FromChannel(ch)
//	result := iter.Filter(func(x int) bool { return x%2 == 0 }).Collect()
func FromChannel[V any](ch <-chan V) *Iterator[V] {
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

// FromChannelWithContext creates a lazy Iterator from a channel with context support.
// The iterator will stop when either the channel is closed or the context is cancelled.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	ch := make(chan int)
//	iter := itertools.FromChannelWithContext(ctx, ch)
func FromChannelWithContext[V any](ctx context.Context, ch <-chan V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// FromReader creates a lazy Iterator that reads lines from an io.Reader.
// Each element is a line from the reader (without the newline character).
// This is useful for processing large files without loading them entirely into memory.
//
// Example:
//
//	file, _ := os.Open("large_file.txt")
//	defer file.Close()
//	iter := itertools.FromReader(file)
//	count := iter.Filter(func(line string) bool {
//	    return strings.Contains(line, "ERROR")
//	}).Count()
func FromReader(r io.Reader) *Iterator[string] {
	return &Iterator[string]{
		seq: func(yield func(string) bool) {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				if !yield(scanner.Text()) {
					return
				}
			}
		},
	}
}

// FromReaderWithContext creates a lazy Iterator from an io.Reader with context support.
// The iterator will stop when either the reader is exhausted or the context is cancelled.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	file, _ := os.Open("large_file.txt")
//	defer file.Close()
//	iter := itertools.FromReaderWithContext(ctx, file)
func FromReaderWithContext(ctx context.Context, r io.Reader) *Iterator[string] {
	return &Iterator[string]{
		seq: func(yield func(string) bool) {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					return
				default:
					if !yield(scanner.Text()) {
						return
					}
				}
			}
		},
	}
}

// FromFunc creates a lazy Iterator from a generator function.
// The function is called repeatedly until it returns false.
// This is useful for generating infinite sequences or custom data sources.
//
// Example:
//
//	// Fibonacci sequence
//	iter := itertools.FromFunc(func() (int, bool) {
//	    a, b := 0, 1
//	    return func() (int, bool) {
//	        result := a
//	        a, b = b, a+b
//	        return result, true
//	    }
//	}())
//	first10 := iter.Take(10).Collect()
func FromFunc[V any](fn func() (V, bool)) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for {
				v, ok := fn()
				if !ok {
					return
				}
				if !yield(v) {
					return
				}
			}
		},
	}
}

// FromFuncWithContext creates a lazy Iterator from a generator function with context support.
// The iterator will stop when either the function returns false or the context is cancelled.
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	iter := itertools.FromFuncWithContext(ctx, func() (int, bool) {
//	    return rand.Int(), true
//	})
//	result := iter.Take(100).Collect()
func FromFuncWithContext[V any](ctx context.Context, fn func() (V, bool)) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					v, ok := fn()
					if !ok {
						return
					}
					if !yield(v) {
						return
					}
				}
			}
		},
	}
}

// Range creates an Iterator that yields integers from start (inclusive) to end (exclusive).
// This is useful for generating sequences of numbers.
//
// Example:
//
//	iter := itertools.Range(0, 10)  // yields 0, 1, 2, ..., 9
//	squares := iter.Map(func(x int) int { return x * x }).Collect()
func Range(start, end int) *Iterator[int] {
	return &Iterator[int]{
		seq: func(yield func(int) bool) {
			for i := start; i < end; i++ {
				if !yield(i) {
					return
				}
			}
		},
	}
}

// RangeStep creates an Iterator that yields integers from start to end with a given step.
//
// Example:
//
//	iter := itertools.RangeStep(0, 10, 2)  // yields 0, 2, 4, 6, 8
//	result := iter.Collect()
func RangeStep(start, end, step int) *Iterator[int] {
	return &Iterator[int]{
		seq: func(yield func(int) bool) {
			if step > 0 {
				for i := start; i < end; i += step {
					if !yield(i) {
						return
					}
				}
			} else if step < 0 {
				for i := start; i > end; i += step {
					if !yield(i) {
						return
					}
				}
			}
		},
	}
}

// Generate creates an infinite Iterator by repeatedly calling a generator function.
// Use Take() or TakeWhile() to limit the output.
//
// Example:
//
//	counter := 0
//	iter := itertools.Generate(func() int {
//	    counter++
//	    return counter
//	})
//	first5 := iter.Take(5).Collect()  // [1, 2, 3, 4, 5]
func Generate[V any](fn func() V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for {
				if !yield(fn()) {
					return
				}
			}
		},
	}
}

// GenerateWithContext creates an infinite Iterator with context support.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//	iter := itertools.GenerateWithContext(ctx, func() int {
//	    return rand.Int()
//	})
//	result := iter.Take(1000).Collect()
func GenerateWithContext[V any](ctx context.Context, fn func() V) *Iterator[V] {
	return &Iterator[V]{
		seq: func(yield func(V) bool) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if !yield(fn()) {
						return
					}
				}
			}
		},
	}
}
