package medium

import "testing"

func TestFindKthLargest(t *testing.T) {
	data := []struct {
		s []int
		k int
		e int
	}{
		{[]int{3, 2, 1, 5, 6, 4}, 2, 5},
		{[]int{3, 2, 3, 1, 2, 4, 5, 5, 6}, 4, 4},
	}
	for i, d := range data {
		got := findKthLargest(d.s, d.k)
		if got != d.e {
			t.Fatalf("%d failed, expect %d got %d", i, d.e, got)
		}
	}
}
