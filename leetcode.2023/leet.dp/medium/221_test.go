package medium_test

import "testing"

func maximalSquare(matrix [][]byte) int {
	m, n := len(matrix), len(matrix[0])
	maxLen := m
	if n < maxLen {
		maxLen = n
	}

	dp := make([][]int16, m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == '1' || matrix[i][j] == 1 {
				dp[i][j] = 1
			}
		}
		for k := 1; k < maxLen; k++ {

		}
	}
	return 0
}

func Test_maximalSquare(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]byte
		want   int
	}{
		{"4-square", [][]byte{{0, 1}, {1, 0}}, 1},
		{"4-square-full", [][]byte{{1, 1}, {1, 1}}, 4},
		{"empty", [][]byte{{0}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maximalSquare(tt.matrix); got != tt.want {
				t.Errorf("maximalSquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
