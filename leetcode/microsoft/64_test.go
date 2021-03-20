package microsoft

import (
	"testing"
)

func minPathSum(grid [][]int) int {
	for i := 1; i < len(grid); i++ {
		grid[i][0] += grid[i-1][0]
	}
	for i := 1; i < len(grid[0]); i++ {
		grid[0][i] += grid[0][i-1]
	}
	for i := 1; i < len(grid); i++ {
		// 已经有好几次写错j < len(grid[0]) 写成 j < len(grid)了，记住啦
		for j := 1; j < len(grid[0]); j++ {
			v1 := grid[i-1][j]
			v2 := grid[i][j-1]
			if v1 < v2 {
				grid[i][j] += v1
			} else {
				grid[i][j] += v2
			}
		}
	}
	return grid[len(grid)-1][len(grid[0])-1]
}

func Test_minPathSum(t *testing.T) {
	type args struct {
		grid [][]int
	}
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{
			grid: [][]int{
				{1, 2, 3},
				{4, 5, 6},
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minPathSum(tt.grid); got != tt.want {
				t.Errorf("minPathSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
