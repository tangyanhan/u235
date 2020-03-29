package easy

import (
	"testing"
)

func TestMergeSortedList(t *testing.T) {
	testData := []struct {
		a      []int
		b      []int
		expect []int
	}{
		{[]int{1, 2, 4}, []int{1, 3, 4}, []int{1, 1, 2, 3, 4, 4}},
	}
	for _, data := range testData {
		a := NewListFromInts(data.a)
		b := NewListFromInts(data.b)
		got := mergeTwoLists(a, b)
		result := ListToInts(got)
		if !CompareInts(result, data.expect) {
			t.Fatalf("a=%v b=%v expect=%v got=%v", data.a, data.b, data.expect, result)
		}
	}
}
