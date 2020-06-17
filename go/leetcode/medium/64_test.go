package medium

import "testing"

func TestFindMin(t *testing.T) {
	testData := []struct {
		grid   [][]int
		expect int
	}{
		{
			grid: [][]int{
				[]int{1, 3, 1},
				[]int{1, 5, 1},
				[]int{4, 2, 1},
			},
			expect: 7,
		},
	}

	for i, data := range testData {
		got := minPathSum(data.grid)
		if got != data.expect {
			t.Fatalf("Data failed at %d: expect %d, got %d", i, data.expect, got)
		}
	}
}
