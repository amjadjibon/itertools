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
