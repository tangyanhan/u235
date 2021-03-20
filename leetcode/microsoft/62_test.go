package microsoft

import "testing"

var memo [100][100]int

func init() {
	for i := 0; i < len(memo); i++ {
		memo[i][0] = 1
		memo[0][i] = 1
	}
}

func uniquePaths(m int, n int) int {
	if v := memo[m-1][n-1]; v != 0 {
		return v
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			// 容易想到i,j的值是通过i-1,j 和i,j-1求和得到,
			// 但最左侧和最上侧都要置为1
			memo[i][j] = memo[i-1][j] + memo[i][j-1]
		}
	}

	return memo[m-1][n-1]
}

func Test_uniquePaths(t *testing.T) {
	tests := []struct {
		name string
		m    int
		n    int
		want int
	}{
		{
			m:    3,
			n:    2,
			want: 3,
		},
		{
			m:    3,
			n:    7,
			want: 28,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniquePaths(tt.m, tt.n); got != tt.want {
				t.Errorf("uniquePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
