package itertools_test

import (
	"fmt"
	"testing"

	"github.com/amjadjibon/itertools"
	"github.com/stretchr/testify/assert"
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

func TestSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	sum := itertools.Sum(iter, func(v int) int { return v }, 0)
	assert.Equal(t, 15, sum)
}

func TestSumFloat(t *testing.T) {
	slice := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	iter := itertools.ToIter(slice)

	sum := itertools.Sum(iter, func(v float64) float64 { return v }, 0)
	assert.Equal(t, 16.5, sum)
}

func TestSumComplex(t *testing.T) {
	type Complex struct {
		A int
		B int
	}

	slice := []Complex{{1, 2}, {3, 4}, {5, 6}}
	iter := itertools.ToIter(slice)

	sum := itertools.Sum(iter, func(v Complex) int { return v.A + v.B }, 0)
	assert.Equal(t, 21, sum)
}

func TestFoldSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	sum := itertools.Fold(iter, func(acc, v int) int { return acc + v }, 0)
	assert.Equal(t, 15, sum)
}

func TestFoldProduct(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	product := itertools.Fold(iter, func(acc, v int) int { return acc * v }, 1)
	assert.Equal(t, 120, product)
}

func TestFoldConcat(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e"}
	iter := itertools.ToIter(slice)

	concat := itertools.Fold(iter, func(acc, v string) string { return acc + v }, "")
	assert.Equal(t, "abcde", concat)
}

func TestProduct(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := itertools.ToIter(slice)

	product := itertools.Product(iter, func(v int) int { return v }, 1)
	assert.Equal(t, 120, product)
}

func TestProductFloat(t *testing.T) {
	slice := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	iter := itertools.ToIter(slice)

	product := itertools.Product(iter, func(v float64) float64 { return v }, 1)
	assert.Equal(t, fmt.Sprintf("%.2f", 1.1*2.2*3.3*4.4*5.5), fmt.Sprintf("%.2f", product))
}

func TestProductComplex(t *testing.T) {
	type Complex struct {
		A int
		B int
	}

	slice := []Complex{{1, 2}, {3, 4}, {5, 6}}
	iter := itertools.ToIter(slice)

	product := itertools.Product(iter, func(v Complex) int { return v.A * v.B }, 1)
	assert.Equal(t, 720, product)
}
