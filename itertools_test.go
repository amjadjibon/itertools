package itertools_test

import (
	"testing"

	"github.com/amjadjibon/itertools"
)

func TestZip(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5, 6}
	slice2 := []string{"a", "b", "c", "d", "e"}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)
	zip := itertools.Zip(iter1, iter2).Collect()

	if len(zip) != 5 {
		t.Errorf("expected 5 elements, got %d", len(zip))
	}

	for i, v := range zip {
		if slice1[i] != v.First {
			t.Errorf("expected %d, got %d", v.First, slice1[i])
		}
		if slice2[i] != v.Second {
			t.Errorf("expected %s, got %s", v.Second, slice2[i])
		}
	}
}

func TestZip2(t *testing.T) {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []string{"a", "b", "c"}
	iter1 := itertools.ToIter(slice1)
	iter2 := itertools.ToIter(slice2)
	zip := itertools.Zip2(iter1, iter2, struct {
		First  int
		Second string
	}{0, ""}).Collect()

	if len(zip) != 5 {
		t.Errorf("expected 5 elements, got %d", len(zip))
	}

	for i, v := range zip {
		if i < len(slice1) {
			if slice1[i] != v.First {
				t.Errorf("expected %d, got %d", v.First, slice1[i])
			}
		} else {
			if v.First != 0 {
				t.Errorf("expected 0, got %d", v.First)
			}
		}

		if i < len(slice2) {
			if slice2[i] != v.Second {
				t.Errorf("expected %s, got %s", v.Second, slice2[i])
			}
		} else {
			if v.Second != "" {
				t.Errorf("expected \"\", got %s", v.Second)
			}
		}
	}
}
