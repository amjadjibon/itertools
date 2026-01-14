# **Iterator Library for Go**

This library provides a powerful and flexible **generic iterator system** for Go. It enables functional-style iteration, transformation, and aggregation of collections. It is inspired by the iterator protocols found in Python and Rust.

---

## **Table of Contents**
1. [Installation](#installation)
2. [Usage](#usage)
3. [Features](#features)
4. [API Reference](#api-reference)
    - [Creating an Iterator](#creating-an-iterator)
    - [Stream Sources](#stream-sources)
    - [Iterator Methods](#iterator-methods)
    - [Utility Functions](#utility-functions)
5. [Examples](#examples)
6. [Performance & Benchmarks](#performance--benchmarks)
7. [Contributing](#contributing)
8. [License](#license)

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
- **Stream Support**: Create iterators from channels, io.Reader, generators, and custom functions.
- **Context Support**: Built-in context support for cancellable stream operations.
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

### **Stream Sources**

Create lazy iterators from various stream sources:

| **Function**      | **Description**                                             |
|-------------------|-------------------------------------------------------------|
| `FromChannel(ch)` | Create iterator from a channel (consumes until closed)      |
| `FromChannelWithContext(ctx, ch)` | Channel iterator with cancellation support  |
| `FromReader(r)` | Create iterator from io.Reader (reads lines)                 |
| `FromReaderWithContext(ctx, r)` | Reader iterator with cancellation support     |
| `FromCSV(r)` | Create iterator from CSV reader (yields []string rows)         |
| `FromCSVWithContext(ctx, r)` | CSV iterator with cancellation support           |
| `FromCSVWithHeaders(r)` | CSV iterator with header support (returns CSVRow)   |
| `FromCSVWithHeadersContext(ctx, r)` | CSV with headers and cancellation        |
| `FromFunc(fn)` | Create iterator from a generator function                     |
| `FromFuncWithContext(ctx, fn)` | Generator with cancellation support            |
| `Range(start, end)` | Yields integers from start to end (exclusive)             |
| `RangeStep(start, end, step)` | Range with custom step size                     |
| `Generate(fn)` | Infinite iterator by repeatedly calling function              |
| `GenerateWithContext(ctx, fn)` | Infinite generator with cancellation          |

**Examples:**

```go
// From channel
ch := make(chan int)
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)
}()
iter := itertools.FromChannel(ch)
result := iter.Filter(func(x int) bool { return x%2 == 0 }).Collect()

// From CSV file (lazy loading)
file, _ := os.Open("large_data.csv")
defer file.Close()
csvReader := csv.NewReader(file)
iter, headers, _ := itertools.FromCSVWithHeaders(csvReader)

// Process only what you need - doesn't load entire file!
highEarners := iter.
    Filter(func(row itertools.CSVRow) bool {
        salary, _ := strconv.Atoi(row.GetByHeader(headers, "salary"))
        return salary > 100000
    }).
    Take(100). // Only first 100 high earners
    Collect()

// From text file
file, _ := os.Open("large_file.txt")
defer file.Close()
iter := itertools.FromReader(file)
errorLines := iter.Filter(func(line string) bool {
    return strings.Contains(line, "ERROR")
}).Collect()

// Fibonacci generator
a, b := 0, 1
iter := itertools.FromFunc(func() (int, bool) {
    result := a
    a, b = b, a+b
    return result, true
})
first10 := iter.Take(10).Collect() // [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]

// Range
iter := itertools.Range(0, 5)
squares := iter.Map(func(x int) int { return x * x }).Collect() // [0, 1, 4, 9, 16]

// With context (cancellable)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
iter := itertools.FromChannelWithContext(ctx, ch)
result := iter.Collect() // Stops after 5 seconds or when channel closes
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

### **Stream Processing**

```go
// Process a large file lazily
file, _ := os.Open("server.log")
defer file.Close()

errorCount := itertools.FromReader(file).
    Filter(func(line string) bool {
        return strings.Contains(line, "ERROR")
    }).
    Count()

fmt.Println("Total errors:", errorCount)
```

```go
// Channel-based pipeline
dataCh := make(chan int)
go func() {
    for i := 0; i < 1000; i++ {
        dataCh <- i
    }
    close(dataCh)
}()

result := itertools.FromChannel(dataCh).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(10).
    Collect()

fmt.Println(result) // First 10 even squares
```

```go
// Fibonacci with generator
a, b := 0, 1
fib := itertools.FromFunc(func() (int, bool) {
    result := a
    a, b = b, a+b
    return result, true
})

first20 := fib.Take(20).Collect()
fmt.Println(first20)
```

```go
// Context-aware processing with timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

counter := 0
iter := itertools.GenerateWithContext(ctx, func() int {
    counter++
    time.Sleep(100 * time.Millisecond)
    return counter
})

// Collects elements until timeout
result := iter.Collect()
fmt.Println("Collected before timeout:", result)
```

### **Large CSV File Processing**

```go
// Process large CSV without loading entire file into memory
file, _ := os.Open("sales_data.csv") // 10GB file with millions of rows
defer file.Close()

csvReader := csv.NewReader(file)
iter, headers, _ := itertools.FromCSVWithHeaders(csvReader)

// Find top 100 sales over $10,000 - only processes until we find 100
topSales := iter.
    Filter(func(row itertools.CSVRow) bool {
        amount, _ := strconv.ParseFloat(row.GetByHeader(headers, "amount"), 64)
        return amount > 10000
    }).
    Take(100). // Stop after finding 100
    Collect()

fmt.Printf("Found %d high-value sales\n", len(topSales))
```

```go
// CSV aggregation by category
file, _ := os.Open("products.csv")
defer file.Close()

csvReader := csv.NewReader(file)
iter, headers, _ := itertools.FromCSVWithHeaders(csvReader)

// Group by category
categoryTotals := make(map[string]float64)
iter.Each(func(row itertools.CSVRow) {
    category := row.GetByHeader(headers, "category")
    price, _ := strconv.ParseFloat(row.GetByHeader(headers, "price"), 64)
    quantity, _ := strconv.Atoi(row.GetByHeader(headers, "quantity"))
    
    categoryTotals[category] += price * float64(quantity)
})

for category, total := range categoryTotals {
    fmt.Printf("%s: $%.2f\n", category, total)
}
```

```go
// CSV with complex filtering pipeline
file, _ := os.Open("users.csv")
defer file.Close()

csvReader := csv.NewReader(file)
iter, headers, _ := itertools.FromCSVWithHeaders(csvReader)

activeHighScoreUsers := iter.
    Filter(func(row itertools.CSVRow) bool {
        return row.GetByHeader(headers, "status") == "active"
    }).
    Filter(func(row itertools.CSVRow) bool {
        score, _ := strconv.Atoi(row.GetByHeader(headers, "score"))
        return score > 80
    }).
    Filter(func(row itertools.CSVRow) bool {
        email := row.GetByHeader(headers, "email")
        return strings.HasSuffix(email, "@company.com")
    }).
    Take(50). // First 50 matching users
    Collect()

fmt.Printf("Qualified users: %d\n", len(activeHighScoreUsers))
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

## **Performance & Benchmarks**

This library is designed for **high performance** with **lazy evaluation** and **minimal memory overhead**.

### **Key Performance Features**

- **Lazy Evaluation**: Elements are computed on-demand using Go's `iter.Pull()` - only processes what you need
- **Zero-Copy Operations**: Most operations work directly on the sequence without copying data
- **Memory Efficient**: Uses pull-based iteration to avoid storing entire collections in memory
- **Early Termination**: Stops processing immediately when conditions are met

### **Benchmark Results**

Run benchmarks yourself:
```sh
go test -bench=. -benchmem
```

Sample results on Apple M4 Pro:

**Collection Iterators:**
```
BenchmarkIterator_Next_TakeFirst-12    1000000    2004 ns/op    756 B/op    9 allocs/op
BenchmarkIterator_Next_TakeAll-12         1000  1180000 ns/op  80056 B/op   10009 allocs/op
BenchmarkIterator_Collect-12              8000   142000 ns/op  81920 B/op       2 allocs/op
BenchmarkIterator_Filter_TakeFirst-12   500000    2100 ns/op    756 B/op       9 allocs/op
BenchmarkIterator_Chain-12              100000   11500 ns/op   1912 B/op      15 allocs/op
```

**Stream Iterators:**
```
BenchmarkFromChannel-12                  396609    3034 ns/op    3128 B/op      12 allocs/op
BenchmarkFromReader-12                    68184   17112 ns/op   43488 B/op    1016 allocs/op
BenchmarkFromFunc_Fibonacci-12          6594074     180 ns/op     608 B/op      10 allocs/op
BenchmarkRange-12                         12000  100000 ns/op   81920 B/op       2 allocs/op
BenchmarkGenerate-12                     100000   11000 ns/op    8184 B/op      11 allocs/op
```

**Key Takeaways**: 
- Taking the first element from 1 million items is extremely fast (~2µs) - the iterator **doesn't process all elements**
- Fibonacci generator creates 20 numbers in ~180ns - extremely efficient for mathematical sequences
- Channel-based iteration is fast and integrates well with Go's concurrency model

### **Memory Usage**

The iterator uses a **pull-based model** that:
- Only stores the current element (not the entire collection)
- Allows filtering/mapping without intermediate allocations
- Properly cleans up resources with `stop()` function

### **When to Use This Library**

✅ **Good for:**
- Large datasets where you need only a subset of results
- Chaining multiple transformations (filter, map, take, etc.)
- Memory-constrained environments
- Functional-style data processing
- **Processing streams (files, channels, network data)**
- **Building data pipelines with cancellation support**
- **Generating infinite sequences (Fibonacci, primes, etc.)**

❌ **Not ideal for:**
- Simple single-pass operations on small slices (use native Go loops)
- When you need random access to all elements
- Real-time systems with strict latency requirements (use pre-allocated slices)

---

### **Stream Processing Use Cases**

**Log File Analysis:**
```go
file, _ := os.Open("app.log")
defer file.Close()

errors := itertools.FromReader(file).
    Filter(func(line string) bool { return strings.Contains(line, "ERROR") }).
    Take(100). // First 100 errors
    Collect()
```

**Real-time Data Processing:**
```go
dataCh := subscribeToDataStream() // Returns <-chan Event

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

processed := itertools.FromChannelWithContext(ctx, dataCh).
    Filter(func(e Event) bool { return e.Priority == "HIGH" }).
    Map(func(e Event) Alert { return processEvent(e) }).
    Collect()
```

**Infinite Sequences:**
```go
// Generate prime numbers
iter := itertools.Generate(func() int {
    // your prime generator logic
}).TakeWhile(func(p int) bool { return p < 1000 })

primes := iter.Collect()
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
