package dp

import "testing"

func TestMoveZero(t *testing.T) {
	d := []int{1, 3, 0, 15, 0, 10}
	moveZeroes(d)
	t.Log(d)
}
