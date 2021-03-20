package microsoft

import "testing"

func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	dp := make([][]int, len(obstacleGrid))
	var colBlock bool
	// 在置1最左边和最上边时，如果遇到阻塞，就要停止置1
	for i := range dp {
		dp[i] = make([]int, len(obstacleGrid[0]))
		if obstacleGrid[i][0] != 0 {
			colBlock = true
		}
		if !colBlock {
			dp[i][0] = 1
		}
	}
	for i := range dp[0] {
		if obstacleGrid[0][i] != 0 {
			break
		}
		dp[0][i] = 1
	}
	for i := 1; i < len(dp); i++ {
		for j := 1; j < len(dp[0]); j++ {
			// 只计算没有阻塞的i,j
			if obstacleGrid[i][j] == 0 {
				dp[i][j] = dp[i-1][j] + dp[i][j-1]
			}
		}
	}
	return dp[len(dp)-1][len(dp[0])-1]
}

func Test_uniquePathsWithObstacles(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		want int
	}{
		{
			grid: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniquePathsWithObstacles(tt.grid); got != tt.want {
				t.Errorf("uniquePathsWithObstacles() = %v, want %v", got, tt.want)
			}
		})
	}
}
