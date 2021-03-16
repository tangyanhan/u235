package microsoft

import (
	"testing"
)

func solution(knows func(a int, b int) bool) func(n int) int {
	return func(n int) int {
		// 在[k, n]区间内，找到任意一个出度为0的候选点
		// 由于名人的性质是入度=n-1，出度=0，当发现任一出度为0的点时，若该点入度不为n-1，则必定没有其它符合条件的点
		var test int
		for i := 1; i < n; i++ {
			if knows(test, i) {
				test = i
			}
		}
		// 检测候选点入度是否为n-1
		for i := 0; i < n; i++ {
			if test == i {
				continue
			}
			if knows(test, i) || !knows(i, test) {
				return -1
			}
		}
		return test
	}
}

func Test_solution(t *testing.T) {
	tests := []struct {
		name  string
		knows func(a int, b int) bool
		n     int
		want  int
	}{
		{
			knows: func(a, b int) bool {
				g := [][]int{{1, 1, 0}, {0, 1, 0}, {1, 1, 1}}
				return g[a][b] == 1
			},
			n:    3,
			want: 1,
		},
		{
			knows: func(a, b int) bool {
				g := [][]int{{1, 0, 0, 0}, {1, 1, 0, 1}, {1, 1, 0, 1}, {0, 0, 0, 1}}
				return g[a][b] == 1
			},
			n:    4,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := solution(tt.knows)
			got := fn(tt.n)
			if got != tt.want {
				t.Errorf("solution() = %v, want %v", got, tt.want)
			}
		})
	}
}
