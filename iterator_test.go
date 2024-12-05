package itertools_test

import (
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

	min, found := iter.Min(func(a, b int) bool {
		return a < b
	})

	assert.True(t, found)
	assert.Equal(t, 1, min)
}

func TestIterator_Min_Empty(t *testing.T) {
	slice := []int{}
	iter := itertools.ToIter(slice)

	min, found := iter.Min(func(a, b int) bool {
		return a < b
	})

	assert.False(t, found)
	assert.Equal(t, 0, min)
}

func TestIterator_Max(t *testing.T) {
	slice := []int{5, 2, 4, 1, 3}
	iter := itertools.ToIter(slice)

	max, found := iter.Max(func(a, b int) bool {
		return a < b
	})

	assert.True(t, found)
	assert.Equal(t, 5, max)
}

func TestIterator_Max_Empty(t *testing.T) {
	slice := []int{}
	iter := itertools.ToIter(slice)

	max, found := iter.Max(func(a, b int) bool {
		return a < b
	})

	assert.False(t, found)
	assert.Equal(t, 0, max)
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
