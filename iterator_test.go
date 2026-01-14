package itertools_test

import (
	"fmt"
	"testing"

	"github.com/amjadjibon/itertools"
	"github.com/stretchr/testify/assert"
)

func TestIterator_Next(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	idx := 0
	for iter.Next() {
		curr := iter.Current()
		if slice[idx] != curr {
			t.Errorf("expected %d, got %d", curr, slice[idx])
		}
		idx++
	}
}

func TestIterator_Next_Empty(t *testing.T) {
	iter := itertools.ToIter([]int{})
	if iter.Next() {
		t.Errorf("expected false, got true")
	}
}

func TestIterator_Current_panics_before_next(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	assert.Panics(t, func() {
		iter.Current()
	})
}

func TestIterator_Current_panics_after_done(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	for iter.Next() {
	}

	assert.Panics(t, func() {
		iter.Current()
	})
}

func TestIterator_Collect(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	collect := iter.Collect()

	assert.Equal(t, slice, collect)
}

func TestIterator_Each(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	collect := make([]int, 0)
	iter.Each(func(v int) {
		collect = append(collect, v)
	})

	assert.Equal(t, slice, collect)
}

func TestIterator_Reverse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	reverse := iter.Reverse().Collect()

	assert.Equal(t, []int{5, 4, 3, 2, 1}, reverse)
}

func TestIterator_Filter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	filtered := iter.Filter(func(v int) bool {
		return v%2 == 0
	}).Collect()

	assert.Equal(t, []int{2, 4}, filtered)
}

func TestIterator_Map(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	mapped := iter.Map(func(v int) int {
		return v * 2
	}).Collect()

	assert.Equal(t, []int{2, 4, 6, 8, 10}, mapped)
}

func TestIterator_Chain(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	chained := iter1.Chain(iter2).Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, chained)
}

func TestIterator_Take(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	taken := iter.Take(3).Collect()

	assert.Equal(t, []int{1, 2, 3}, taken)
}

func TestIterator_Drop(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	dropped := iter.Drop(2).Collect()

	assert.Equal(t, []int{3, 4, 5}, dropped)
}

func TestIterator_TakeWhile(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	taken := iter.TakeWhile(func(v int) bool {
		return v < 4
	}).Collect()

	assert.Equal(t, []int{1, 2, 3}, taken)
}

func TestIterator_DropWhile(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	dropped := iter.DropWhile(func(v int) bool {
		return v < 3
	}).Collect()

	assert.Equal(t, []int{3, 4, 5}, dropped)
}

func TestIterator_First(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	first := iter.First()

	assert.Equal(t, 1, first)
}

func TestIterator_Last(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	last := iter.Last()

	assert.Equal(t, 5, last)
}

func TestIterator_Nth(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	nth := iter.Nth(2)

	assert.Equal(t, 3, nth)
}

func TestIterator_Any(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	anyResult := iter.Any(func(v int) bool {
		return v == 3
	})

	assert.True(t, anyResult)
}

func TestIterator_All(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	all := iter.All(func(v int) bool {
		return v < 6
	})

	assert.True(t, all)
}

func TestIterator_Find(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	result, found := iter.Find(func(v int) bool {
		return v == 3
	})

	assert.True(t, found)
	assert.Equal(t, 3, result)
}

func TestIterator_Find_NotFound(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	result, found := iter.Find(func(v int) bool {
		return v == 6
	})

	assert.False(t, found)
	assert.Equal(t, 0, result)
}

func TestIterator_Sort(t *testing.T) {
	slice := []int{5, 2, 4, 1, 3}
	iter := itertools.ToIter(slice)

	sorted := iter.Sort(func(a, b int) bool {
		return a < b
	}).Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5}, sorted)
}

func TestIterator_Min(t *testing.T) {
	slice := []int{5, 2, 4, 1, 3}
	iter := itertools.ToIter(slice)

	minValue, found := iter.Min(func(a, b int) bool {
		return a < b
	})

	assert.True(t, found)
	assert.Equal(t, 1, minValue)
}

func TestIterator_Min_Empty(t *testing.T) {
	slice := []int{}
	iter := itertools.ToIter(slice)

	minValue, found := iter.Min(func(a, b int) bool {
		return a < b
	})

	assert.False(t, found)
	assert.Equal(t, 0, minValue)
}

func TestIterator_Max(t *testing.T) {
	slice := []int{5, 2, 4, 1, 3}
	iter := itertools.ToIter(slice)

	maxValue, found := iter.Max(func(a, b int) bool {
		return a < b
	})

	assert.True(t, found)
	assert.Equal(t, 5, maxValue)
}

func TestIterator_Max_Empty(t *testing.T) {
	slice := []int{}
	iter := itertools.ToIter(slice)

	maxValue, found := iter.Max(func(a, b int) bool {
		return a < b
	})

	assert.False(t, found)
	assert.Equal(t, 0, maxValue)
}

func TestIterator_Partition(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	even, odd := iter.Partition(func(v int) bool {
		return v%2 == 0
	})

	assert.Equal(t, []int{2, 4}, even.Collect())
	assert.Equal(t, []int{1, 3, 5}, odd.Collect())
}

func TestIterator_Count(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	count := iter.Count()

	assert.Equal(t, 5, count)
}

func TestIterator_Unique(t *testing.T) {
	slice := []int{1, 2, 2, 3, 3, 4, 5}
	iter := itertools.ToIter(slice)

	unique := iter.Unique(func(i int) any {
		return i
	})

	assert.Equal(t, []int{1, 2, 3, 4, 5}, unique.Collect())
}

func TestIterator_Unique_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Alice", 25},
		{"Charlie", 35},
		{"Bob", 30},
	}

	iter := itertools.ToIter(slice)

	unique := iter.Unique(func(p person) any {
		return p.Name
	})

	assert.Equal(t, []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}, unique.Collect())
}

func TestIterator_GroupBy(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Alice", 26},
		{"Charlie", 35},
		{"Bob", 43},
	}

	iter := itertools.ToIter(slice)

	groups := iter.GroupBy(func(p person) string {
		return p.Name
	})

	assert.Equal(t, map[string][]person{
		"Alice":   {{Name: "Alice", Age: 25}, {Name: "Alice", Age: 26}},
		"Bob":     {{Name: "Bob", Age: 30}, {Name: "Bob", Age: 43}},
		"Charlie": {{Name: "Charlie", Age: 35}},
	}, groups)

	alices := itertools.ToIter(groups["Alice"])
	alicesAges := itertools.Sum(alices, func(p person) int { return p.Age }, 0)
	assert.Equal(t, 51, alicesAges)

	bobs := itertools.ToIter(groups["Bob"])
	bobsMaxAge, _ := bobs.Max(func(a, b person) bool { return a.Age < b.Age })
	assert.Equal(t, person{Name: "Bob", Age: 43}, bobsMaxAge)

	charlies := itertools.ToIter(groups["Charlie"])
	charliesCount := charlies.Count()
	assert.Equal(t, 1, charliesCount)
}

func TestIterator_Cycle(t *testing.T) {
	slice := []int{1, 2, 3}
	iter := itertools.ToIter(slice)

	cycle := iter.Cycle().Take(7).Collect()

	assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1}, cycle)
}

func TestIterator_Repeat(t *testing.T) {
	iter := itertools.Repeat(42, 3)

	repeat := iter.Collect()

	assert.Equal(t, []int{42, 42, 42}, repeat)
}

func TestIterator_Union(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{3, 4, 5}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	union := iter1.Union(iter2, func(i int) any {
		return i
	}).Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5}, union)
}

func TestIterator_Union_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice1 := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}
	slice2 := []person{
		{"Charlie", 35},
		{"David", 40},
		{"Alice", 25},
	}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	union := iter1.Union(iter2, func(p person) any {
		return p.Name
	}).Collect()

	assert.Equal(t, []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
		{"David", 40},
	}, union)
}

func TestIterator_Intersection(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{3, 4, 5}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	intersection := iter1.Intersection(iter2, func(i int) any {
		return i
	}).Collect()

	assert.Equal(t, []int{3}, intersection)
}

func TestIterator_Intersection_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice1 := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}
	slice2 := []person{
		{"Charlie", 35},
		{"David", 40},
		{"Alice", 25},
	}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	intersection := iter1.Intersection(iter2, func(p person) any {
		return p.Name
	}).Collect()

	assert.Equal(t, []person{
		{"Alice", 25},
		{"Charlie", 35},
	}, intersection)
}

func TestIterator_Difference(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{3, 4, 5}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	difference := iter1.Difference(iter2, func(i int) any {
		return i
	}).Collect()

	assert.Equal(t, []int{1, 2}, difference)
}

func TestIterator_Difference_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice1 := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}
	slice2 := []person{
		{"Charlie", 35},
		{"David", 40},
		{"Alice", 25},
	}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)

	difference := iter1.Difference(iter2, func(p person) any {
		return p.Name
	}).Collect()

	assert.Equal(t, []person{
		{"Bob", 30},
	}, difference)
}

func TestIterator_StepBy(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	iter := itertools.ToIter(slice)

	stepBy := iter.StepBy(3).Collect()

	assert.Equal(t, []int{1, 4, 7}, stepBy)
}

func TestIterator_Shuffle(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	shuffled := iter.Shuffle()

	once := shuffled.Collect()
	assert.ElementsMatch(t, slice, once)
	assert.NotEqual(t, slice, once)
}

func TestIterator_Index(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	index := iter.Index(func(v int) bool {
		return v == 3
	})

	assert.Equal(t, 2, index)
}

func TestIterator_Index_NotFound(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	index := iter.Index(func(v int) bool {
		return v == 6
	})

	assert.Equal(t, -1, index)
}

func TestIterator_Index_Complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}
	iter := itertools.ToIter(slice)

	index := iter.Index(func(p person) bool {
		return p.Name == "Bob"
	})

	assert.Equal(t, 1, index)
}

func TestIterator_LastIndex(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 3}
	iter := itertools.ToIter(slice)

	index := iter.LastIndex(func(v int) bool {
		return v == 3
	})

	assert.Equal(t, 5, index)
}

func TestIterator_IsSorted(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	sorted := iter.IsSorted(func(a, b int) bool {
		return a < b
	})

	assert.True(t, sorted)
}

func TestIterator_String(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	str := iter.String()

	assert.Equal(t, fmt.Sprintf("<Iterator: %v>", slice), str)
}

func TestIterator_Replace(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	replaced := iter.Replace(func(v int) bool { return v%2 == 0 }, 0).Collect()

	assert.Equal(t, []int{1, 0, 3, 0, 5}, replaced)
}

func TestIterator_ReplaceAll(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	replaced := iter.ReplaceAll(0).Collect()

	assert.Equal(t, []int{0, 0, 0, 0, 0}, replaced)
}

func TestIterator_Compact(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 0, 6, 0, 7}
	iter := itertools.ToIter(slice)

	compacted := iter.Compact().Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, compacted)
}

func TestIterator_Compact_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice := []*person{
		{"Alice", 25},
		nil,
		{"Charlie", 35},
	}
	iter := itertools.ToIter(slice)

	compacted := iter.Compact().Collect()

	assert.Equal(t, []*person{
		{"Alice", 25},
		{"Charlie", 35},
	}, compacted)
}

func TestIterator_CompactWith(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 0, 6, 0, 7}
	iter := itertools.ToIter(slice)

	compacted := iter.CompactWith(0).Collect()

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, compacted)
}

func TestIterator_CompactWith_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice := []*person{
		{"Alice", 25},
		{Name: "", Age: 0},
		{"Charlie", 35},
	}
	iter := itertools.ToIter(slice)

	compacted := iter.CompactWith(&person{}).Collect()

	assert.Equal(t, []*person{
		{"Alice", 25},
		{"Charlie", 35},
	}, compacted)
}

func TestIterator_ToUpper(t *testing.T) {
	slice := []string{"hello", "world"}
	iter := itertools.ToIter(slice)

	upper := iter.ToUpper().Collect()

	assert.Equal(t, []string{"HELLO", "WORLD"}, upper)

	slice2 := []int{1, 2, 3}
	iter2 := itertools.ToIter(slice2)
	upper2 := iter2.ToUpper().Collect()

	assert.Equal(t, []int{1, 2, 3}, upper2)
}

func TestIterator_ToLower(t *testing.T) {
	slice := []string{"HELLO", "WORLD"}
	iter := itertools.ToIter(slice)

	lower := iter.ToLower().Collect()

	assert.Equal(t, []string{"hello", "world"}, lower)

	slice2 := []int{1, 2, 3}
	iter2 := itertools.ToIter(slice2)
	lower2 := iter2.ToLower().Collect()

	assert.Equal(t, []int{1, 2, 3}, lower2)
}

func TestIterator_TrimSpace(t *testing.T) {
	slice := []string{"  hello  ", "world  "}
	iter := itertools.ToIter(slice)

	trimmed := iter.TrimSpace().Collect()

	assert.Equal(t, []string{"hello", "world"}, trimmed)
}

// AssertEq is a helper function to compare two slices of any type.
func TestIterator_AssertEq(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	res := iter.AssertEq(slice, func(i1, i2 int) bool {
		return i1 == i2
	})

	assert.True(t, res)
}

func TestIterator_AssertEq_complex(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	slice := []person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}
	iter := itertools.ToIter(slice)

	res := iter.AssertEq(slice, func(p1, p2 person) bool {
		return p1.Name == p2.Name && p1.Age == p2.Age
	})

	assert.True(t, res)
}

// TestChainEarlyTermination tests that Chain stops the second iterator when yield returns false
func TestChainEarlyTermination(t *testing.T) {
	first := itertools.ToIter([]int{1, 2, 3})
	second := itertools.ToIter([]int{4, 5, 6})

	chained := first.Chain(second)

	// Take only 2 elements - should stop before reaching second iterator
	result := chained.Take(2).Collect()

	expected := []int{1, 2}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
		}
	}
}

// TestChainFullIteration tests that Chain works correctly when fully consumed
func TestChainFullIteration(t *testing.T) {
	first := itertools.ToIter([]int{1, 2, 3})
	second := itertools.ToIter([]int{4, 5, 6})

	chained := first.Chain(second)
	result := chained.Collect()

	expected := []int{1, 2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, result[i])
		}
	}
}

// TestUniqueMultipleIterations tests that Unique works correctly when iterated multiple times
func TestUniqueMultipleIterations(t *testing.T) {
	data := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	iter := itertools.ToIter(data).Unique(func(v int) any { return v })

	// First iteration
	first := iter.Collect()
	expected := []int{1, 2, 3, 4}

	if len(first) != len(expected) {
		t.Errorf("First iteration: Expected length %d, got %d", len(expected), len(first))
	}

	for i, v := range expected {
		if first[i] != v {
			t.Errorf("First iteration: Expected %d at index %d, got %d", v, i, first[i])
		}
	}

	// Second iteration - should produce same results (fresh map)
	iter2 := itertools.ToIter(data).Unique(func(v int) any { return v })
	second := iter2.Collect()

	if len(second) != len(expected) {
		t.Errorf("Second iteration: Expected length %d, got %d", len(expected), len(second))
	}

	for i, v := range expected {
		if second[i] != v {
			t.Errorf("Second iteration: Expected %d at index %d, got %d", v, i, second[i])
		}
	}
}

// TestFirstOrWithElements tests FirstOr with non-empty iterator
func TestFirstOrWithElements(t *testing.T) {
	iter := itertools.ToIter([]int{1, 2, 3})
	result := iter.FirstOr(999)

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}

// TestFirstOrWithEmptyIterator tests FirstOr with empty iterator
func TestFirstOrWithEmptyIterator(t *testing.T) {
	iter := itertools.ToIter([]int{})
	result := iter.FirstOr(999)

	if result != 999 {
		t.Errorf("Expected 999, got %d", result)
	}
}

// TestLastOrWithElements tests LastOr with non-empty iterator
func TestLastOrWithElements(t *testing.T) {
	iter := itertools.ToIter([]int{1, 2, 3})
	result := iter.LastOr(999)

	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
}

// TestLastOrWithEmptyIterator tests LastOr with empty iterator
func TestLastOrWithEmptyIterator(t *testing.T) {
	iter := itertools.ToIter([]int{})
	result := iter.LastOr(999)

	if result != 999 {
		t.Errorf("Expected 999, got %d", result)
	}
}

// TestNthOrWithElements tests NthOr with valid index
func TestNthOrWithElements(t *testing.T) {
	iter := itertools.ToIter([]int{10, 20, 30, 40})
	result := iter.NthOr(2, 999)

	if result != 30 {
		t.Errorf("Expected 30, got %d", result)
	}
}

// TestNthOrWithInvalidIndex tests NthOr with index out of bounds
func TestNthOrWithInvalidIndex(t *testing.T) {
	iter := itertools.ToIter([]int{10, 20, 30})
	result := iter.NthOr(5, 999)

	if result != 999 {
		t.Errorf("Expected 999, got %d", result)
	}
}

// TestNthOrWithEmptyIterator tests NthOr with empty iterator
func TestNthOrWithEmptyIterator(t *testing.T) {
	iter := itertools.ToIter([]int{})
	result := iter.NthOr(0, 999)

	if result != 999 {
		t.Errorf("Expected 999, got %d", result)
	}
}

// TestIsSortedWithEmptyIterator tests IsSorted with empty iterator
func TestIsSortedWithEmptyIterator(t *testing.T) {
	iter := itertools.ToIter([]int{})
	result := iter.IsSorted(func(a, b int) bool { return a < b })

	if !result {
		t.Errorf("Expected true for empty iterator, got false")
	}
}

// TestIsSortedWithSingleElement tests IsSorted with single element
func TestIsSortedWithSingleElement(t *testing.T) {
	iter := itertools.ToIter([]int{42})
	result := iter.IsSorted(func(a, b int) bool { return a < b })

	if !result {
		t.Errorf("Expected true for single element, got false")
	}
}

// TestIsSortedWithSortedElements tests IsSorted with sorted elements
func TestIsSortedWithSortedElements(t *testing.T) {
	iter := itertools.ToIter([]int{1, 2, 3, 4, 5})
	result := iter.IsSorted(func(a, b int) bool { return a < b })

	if !result {
		t.Errorf("Expected true for sorted elements, got false")
	}
}

// TestIsSortedWithUnsortedElements tests IsSorted with unsorted elements
func TestIsSortedWithUnsortedElements(t *testing.T) {
	iter := itertools.ToIter([]int{1, 3, 2, 4, 5})
	result := iter.IsSorted(func(a, b int) bool { return a < b })

	if result {
		t.Errorf("Expected false for unsorted elements, got true")
	}
}

// TestIsSortedWithDuplicates tests IsSorted with duplicate elements
func TestIsSortedWithDuplicates(t *testing.T) {
	iter := itertools.ToIter([]int{1, 2, 2, 3, 4})
	result := iter.IsSorted(func(a, b int) bool { return a < b })

	if !result {
		t.Errorf("Expected true for sorted elements with duplicates, got false")
	}
}

// TestUniqueWithComplexKeys tests Unique with complex key function
func TestUniqueWithComplexKeys(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	data := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Alice", 35}, // Duplicate name
		{"Charlie", 25},
		{"Bob", 25}, // Duplicate name
	}

	iter := itertools.ToIter(data).Unique(func(p Person) any { return p.Name })
	result := iter.Collect()

	expected := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 25},
	}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i].Name != v.Name {
			t.Errorf("Expected name %s at index %d, got %s", v.Name, i, result[i].Name)
		}
	}
}
