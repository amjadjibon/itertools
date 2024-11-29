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
