package offer100_test

import "testing"

func maxValue(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	for col := 0; col < n; col++ {
		for row := 0; row < m; row++ {
			left := col - 1
			up := row - 1
			var leftVal, upVal int
			if left >= 0 {
				leftVal = grid[row][left]
			}
			if up >= 0 {
				upVal = grid[up][col]
			}

			if leftVal > upVal {
				grid[row][col] += leftVal
			} else {
				grid[row][col] += upVal
			}
		}
	}

	return grid[m-1][n-1]
}

func Test_maxValue(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{"general", [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}, 12},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxValue(tt.grid); got != tt.want {
				t.Errorf("maxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
