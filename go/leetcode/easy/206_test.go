package easy

import "testing"

func TestReverseList(t *testing.T) {
	testData := []struct {
		a      []int
		expect []int
	}{
		{[]int{1, 2, 4}, []int{4, 2, 1}},
		{[]int{}, []int{}},
	}
	for _, data := range testData {
		a := NewListFromInts(data.a)
		got := reverseList(a)
		result := ListToInts(got)
		if !CompareInts(result, data.expect) {
			t.Fatalf("a=%v expect=%v got=%v", data.a, data.expect, result)
		}
	}
}
