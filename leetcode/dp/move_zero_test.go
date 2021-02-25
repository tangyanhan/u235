package dp

import "testing"

func TestMoveZero(t *testing.T) {
	d := []int{1, 3, 0, 15, 0, 10}
	moveZeroes(d)
	t.Log(d)
	d = []int{0, 1, 0, 3, 12}
	moveZeroes(d)
	t.Log(d)
}
