# `itertools`

`itertools` is a lightweight, generic Go library that provides powerful and reusable iterator utilities. It allows you to iterate over sequences, apply transformations, filter data, and compose operations in a functional programming style.

---

## Features

- **Generic Iterators**: Iterate over any type of sequence.
- **Transformations**: Map, filter, and chain iterators.
- **Collection Utilities**: Collect elements into slices, apply functions to each element, and more.
- **Control Operations**: Take or drop elements based on conditions or counts.
- **Reversals**: Iterate over sequences in reverse order.

---

## Installation

To use `itertools` in your Go project, simply install the library:

```bash
go get github.com/amjadjibon/itertools
```

---

## Usage

### Create an Iterator

You can create an iterator from a slice using `ToIter`:

```go
package main

import (
	"fmt"
    
	"github.com/amjadjibon/itertools"
)

func main() {
	data := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(data)

	for iter.Next() {
		fmt.Println(iter.Current())
	}
}
```

### Transformations

#### Map

Apply a function to each element of the iterator:

```go
mapped := iter.Map(func(v int) int {
	return v * 2
})
fmt.Println(mapped.Collect()) // [2, 4, 6, 8, 10]
```

#### Filter

Filter elements based on a predicate:

```go
filtered := iter.Filter(func(v int) bool {
	return v%2 == 0
})
fmt.Println(filtered.Collect()) // [2, 4]
```

#### Reverse

Reverse the sequence:

```go
reversed := iter.Reverse()
fmt.Println(reversed.Collect()) // [5, 4, 3, 2, 1]
```

### Control Operations

#### Take

Take the first `n` elements:

```go
taken := iter.Take(3)
fmt.Println(taken.Collect()) // [1, 2, 3]
```

#### Drop

Skip the first `n` elements:

```go
dropped := iter.Drop(3)
fmt.Println(dropped.Collect()) // [4, 5]
```

#### TakeWhile

Take elements while a predicate is true:

```go
takeWhile := iter.TakeWhile(func(v int) bool {
	return v < 4
})
fmt.Println(takeWhile.Collect()) // [1, 2, 3]
```

#### DropWhile

Skip elements while a predicate is true:

```go
dropWhile := iter.DropWhile(func(v int) bool {
	return v < 4
})
fmt.Println(dropWhile.Collect()) // [4, 5]
```

### Combining Iterators

#### Chain

Concatenate two iterators:

```go
otherData := []int{6, 7, 8}
otherIter := itertools.ToIter(otherData)

chained := iter.Chain(otherIter)
fmt.Println(chained.Collect()) // [1, 2, 3, 4, 5, 6, 7, 8]
```

---

## API Reference

### Methods

#### `Next() bool`
Advances the iterator to the next element. Returns `true` if an element is available.

#### `Current() V`
Returns the current element. Panics if the iterator has not started or is done.

#### `Collect() []V`
Collects all elements from the iterator into a slice.

#### `Each(f func(V))`
Applies a function to each element of the iterator.

#### `Reverse() *Iterator[V]`
Returns an iterator that iterates over elements in reverse order.

#### `Filter(predicate func(V) bool) *Iterator[V]`
Yields only elements that satisfy the predicate.

#### `Map(f func(V) V) *Iterator[V]`
Transforms each element using the provided function.

#### `Chain(other *Iterator[V]) *Iterator[V]`
Concatenates two iterators.

#### `Take(n int) *Iterator[V]`
Yields the first `n` elements.

#### `Drop(n int) *Iterator[V]`
Skips the first `n` elements.

#### `TakeWhile(predicate func(V) bool) *Iterator[V]`
Yields elements while the predicate is true.

#### `DropWhile(predicate func(V) bool) *Iterator[V]`
Skips elements while the predicate is true.

---

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

---

## License

This library is open-source and available under the [MIT License](LICENSE).

---

Happy iterating! 🎉