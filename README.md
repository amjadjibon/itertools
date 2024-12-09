# **Iterator Library for Go**

This library provides a powerful and flexible **generic iterator system** for Go. It enables functional-style iteration, transformation, and aggregation of collections. It is inspired by the iterator protocols found in Python and Rust.

---

## **Table of Contents**
1. [Installation](#installation)
2. [Usage](#usage)
3. [Features](#features)
4. [API Reference](#api-reference)
    - [Creating an Iterator](#creating-an-iterator)
    - [Iterator Methods](#iterator-methods)
    - [Utility Functions](#utility-functions)
5. [Examples](#examples)
6. [Contributing](#contributing)
7. [License](#license)

---

## **Installation**

```sh
go get -u github.com/amjadjibon/itertools
```

---

## **Usage**

Here's a simple example to get you started.

```go
package main

import (
    "fmt"
    "github.com/amjadjibon/itertools"
)

func main() {
    // Create an iterator from a slice
    iter := itertools.ToIter([]int{1, 2, 3, 4, 5})

    // Filter even numbers, map them to their squares, and collect the result
    result := iter.Filter(func(x int) bool { return x%2 == 0 }).
        Map(func(x int) int { return x * x }).
        Collect()

    fmt.Println(result) // Output: [4, 16]
}
```

---

## **Features**
- **Chainable API**: Combine transformations like `Filter`, `Map`, `Take`, and `Drop` into one functional-style chain.
- **Laziness**: Iterators are lazy; they only compute elements as needed.
- **Composable**: Supports operations like `Zip`, `Chain`, `Union`, `Intersection`, `Difference`, and `Flatten`.
- **Collection Methods**: Collect items into slices, count them, partition them, and more.
- **Generalized Iterators**: Supports all types, as it uses Go's generics.
- **Multiple Data Transformations**: Sort, shuffle, reverse, compact, and manipulate iterator contents.

---

## **API Reference**

### **Creating an Iterator**

1. **From Slice**
    ```go
    iter := itertools.ToIter([]int{1, 2, 3, 4})
    ```

2. **Custom Sequences**
    ```go
    iter := itertools.NewIterator(1, 2, 3, 4)
    ```

3. **Repeated Values**
    ```go
    iter := itertools.Repeat(42, 5)
    ```

4. **Cycle**
    ```go
    iter := itertools.NewIterator(1, 2, 3).Cycle()
    ```

---

### **Iterator Methods**

These methods modify or operate on the elements of an iterator.

| **Method**       | **Description**                                             |
|------------------|------------------------------------------------------------|
| `Next()`         | Advances the iterator to the next element.                   |
| `Current()`      | Returns the current element.                                 |
| `Collect()`      | Collects all elements into a slice.                          |
| `Each(f func(V))`| Applies `f` to each element.                                 |
| `Filter(f func(V) bool)` | Yields only elements that satisfy the predicate `f`.|
| `Map(f func(V) V)` | Transforms each element using `f`.                         |
| `Reverse()`      | Iterates over elements in reverse order.                     |
| `Take(n int)`    | Takes the first `n` elements.                                |
| `Drop(n int)`    | Skips the first `n` elements.                                |
| `TakeWhile(f func(V) bool)` | Yields elements while the predicate `f` is true.|
| `DropWhile(f func(V) bool)` | Drops elements while the predicate `f` is true.  |
| `Count()`        | Counts the total number of elements.                         |
| `First()`        | Returns the first element.                                   |
| `Last()`         | Returns the last element.                                    |
| `Nth(n int)`     | Returns the nth element.                                     |
| `Sort(less func(a, b V) bool)` | Sorts elements according to `less`.          |
| `Min(less func(a, b V) bool)` | Returns the minimum element.                   |
| `Max(less func(a, b V) bool)` | Returns the maximum element.                   |
| `Any(f func(V) bool)` | Returns true if any element satisfies `f`.             |
| `All(f func(V) bool)` | Returns true if all elements satisfy `f`.              |
| `Find(f func(V) bool)` | Returns the first element that satisfies `f`.         |
| `Index(f func(V) bool)` | Returns the index of the first element that satisfies `f`. |
| `LastIndex(f func(V) bool)` | Returns the index of the last element that satisfies `f`. |
| `IsSorted(less func(a, b V) bool)` | Returns true if the elements are sorted.  |
| `Replace(f func(V) bool, replacement V)` | Replaces elements that satisfy `f`.|
| `Compact()`      | Removes nil/zero-value elements.                             |
| `Union(other *Iterator, keyFunc func(V) any)` | Merges two iterators without duplicates.|
| `Difference(other *Iterator, keyFunc func(V) any)` | Difference of two iterators.|
| `Intersection(other *Iterator, keyFunc func(V) any)` | Intersection of two iterators.|

---

### **Utility Functions**

| **Function**      | **Description**                                             |
|-------------------|------------------------------------------------------------|
| `Zip(it1, it2)`   | Zips two iterators together.                                 |
| `Zip2(it1, it2, fill)` | Zips two iterators, filling extra elements with `fill`.|
| `Fold(it, transform, initial)` | Reduces the elements using `transform`.       |
| `Sum(it, transform, zero)` | Sums the elements.                                 |
| `Product(it, transform, one)` | Computes the product of elements.             |
| `ChunkSlice(it, size)` | Returns slices of `size`.                              |
| `Flatten(it1, it2, ...)` | Flattens multiple iterators into one.                |
| `CartesianProduct(it1, it2)` | Generates Cartesian product of two iterators.  |

---

## **Examples**

### **Basic Usage**

```go
// Create an iterator from a slice
iter := itertools.ToIter([]int{1, 2, 3, 4})

// Filter and Map
result := iter.Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * 2 }).
    Collect()

fmt.Println(result) // Output: [4, 8]
```

### **Sort Elements**

```go
iter := itertools.ToIter([]int{3, 1, 4, 2})
sorted := iter.Sort(func(a, b int) bool { return a < b }).Collect()
fmt.Println(sorted) // Output: [1, 2, 3, 4]
```

### **Zip Two Iterators**

```go
iter1 := itertools.ToIter([]int{1, 2, 3})
iter2 := itertools.ToIter([]string{"a", "b", "c"})

zipped := itertools.Zip(iter1, iter2).Collect()
fmt.Println(zipped) 
// Output: [{1 a} {2 b} {3 c}]
```

### **Generate Cartesian Product**

```go
iter1 := itertools.ToIter([]int{1, 2})
iter2 := itertools.ToIter([]string{"a", "b"})

cartesian := itertools.CartesianProduct(iter1, iter2).Collect()
fmt.Println(cartesian)
// Output: [{X: 1, Y: "a"} {X: 1, Y: "b"} {X: 2, Y: "a"} {X: 2, Y: "b"}]
```

---

## **Contributing**

Contributions are welcome! To get started:
1. Fork the repository.
2. Create a new branch for your feature/bugfix.
3. Submit a pull request.

---

## **License**

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

---

With this library, you can process collections in a functional, chainable, and lazy manner. From filtering and mapping to complex operations like cartesian products, this iterator system brings the power of iterables to Go. Happy iterating! 🎉
