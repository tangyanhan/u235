package medium

import "testing"

func TestRemoveNthFromEnd(t *testing.T) {
	s := []struct {
		n int
		s []int
		e []int
	}{
		{
			n: 2,
			s: []int{1, 2, 3, 4, 5},
			e: []int{1, 2, 3, 5},
		},
		{
			n: 5,
			s: []int{1, 2, 3, 4, 5},
			e: []int{2, 3, 4, 5},
		},
		{
			n: 1,
			s: []int{1},
			e: []int{},
		},
	}
	for i, v := range s {
		l := CreateList(v.s)
		l2 := removeNthFromEnd(l, v.n)
		got := ListToSlice(l2)
		if !CmpSlice(v.e, got) {
			t.Fatalf("%d failed, expect %v got %v", i, v.e, got)
		}
	}
}
