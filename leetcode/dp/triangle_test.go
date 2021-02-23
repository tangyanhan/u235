package dp

import "testing"

func TestTriangle(t *testing.T) {
	d := [][]int{{2}, {3, 4}, {6, 5, 7}, {4, 1, 8, 3}}
	v := minimumTotal(d)
	if v != 11 {
		t.Fatal(v)
	}
}
