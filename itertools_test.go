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

func TestIterator_Chain(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)
	chain := iter1.Chain(iter2).Collect()

	if len(chain) != 6 {
		t.Errorf("expected 6 elements, got %d", len(chain))
	}

	for i, v := range chain {
		if i < 3 {
			if slice1[i] != v {
				t.Errorf("expected %d, got %d", v, slice1[i])
			}
		} else {
			if slice2[i-3] != v {
				t.Errorf("expected %d, got %d", v, slice2[i-3])
			}
		}
	}
}

func TestIterator_Take(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	taken := iter.Take(3).Collect()

	if len(taken) != 3 {
		t.Errorf("expected 3 elements, got %d", len(taken))
	}

	for i, v := range taken {
		if slice[i] != v {
			t.Errorf("expected %d, got %d", v, slice[i])
		}
	}
}

func TestIterator_Drop(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	dropped := iter.Drop(2).Collect()

	if len(dropped) != 3 {
		t.Errorf("expected 3 elements, got %d", len(dropped))
	}

	for i, v := range dropped {
		if slice[i+2] != v {
			t.Errorf("expected %d, got %d", v, slice[i+2])
		}
	}
}

func TestIterator_TakeWhile(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	taken := iter.TakeWhile(func(v int) bool {
		return v < 4
	}).Collect()

	if len(taken) != 3 {
		t.Errorf("expected 3 elements, got %d", len(taken))
	}

	for i, v := range taken {
		if slice[i] != v {
			t.Errorf("expected %d, got %d", v, slice[i])
		}
	}
}

func TestIterator_DropWhile(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)
	dropped := iter.DropWhile(func(v int) bool {
		return v < 3
	}).Collect()

	if len(dropped) != 3 {
		t.Errorf("expected 3 elements, got %d", len(dropped))
	}

	for i, v := range dropped {
		if slice[i+2] != v {
			t.Errorf("expected %d, got %d", v, slice[i+2])
		}
	}
}
