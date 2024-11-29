package itertools_test

import (
	"testing"

	"github.com/amjadjibon/itertools"
)

func TestIterator_Collect(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	collect := iter.Collect()

	if len(collect) != 5 {
		t.Errorf("expected 5 elements, got %d", len(collect))
	}

	for i, v := range collect {
		if slice[i] != v {
			t.Errorf("expected %d, got %d", v, slice[i])
		}
	}
}

func TestIterator_Each(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	collect := make([]int, 0)
	iter.Each(func(v int) {
		collect = append(collect, v)
	})

	if len(collect) != 5 {
		t.Errorf("expected 5 elements, got %d", len(collect))
	}

	for i, v := range collect {
		if slice[i] != v {
			t.Errorf("expected %d, got %d", v, slice[i])
		}
	}
}

func TestIterator_Reverse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	reverse := iter.Reverse().Collect()

	if len(reverse) != 5 {
		t.Errorf("expected 5 elements, got %d", len(reverse))
	}

	for i, v := range reverse {
		if slice[4-i] != v {
			t.Errorf("expected %d, got %d", v, slice[4-i])
		}
	}
}

func TestIterator_Filter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	filtered := iter.Filter(func(v int) bool {
		return v%2 == 0
	}).Collect()

	if len(filtered) != 2 {
		t.Errorf("expected 2 elements, got %d", len(filtered))
	}

	for i, v := range filtered {
		if slice[1+i*2] != v {
			t.Errorf("expected %d, got %d", v, slice[1+i*2])
		}
	}
}

func TestIterator_Map(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	mapped := iter.Map(func(v int) int {
		return v * 2
	}).Collect()

	if len(mapped) != 5 {
		t.Errorf("expected 5 elements, got %d", len(mapped))
	}

	for i, v := range mapped {
		if slice[i]*2 != v {
			t.Errorf("expected %d, got %d", v, slice[i]*2)
		}
	}
}
