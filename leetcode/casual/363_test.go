package casual

import (
	"math"
	"testing"
)

func maxSumSubmatrix(matrix [][]int, k int) int {
	maxValue := math.MinInt32
	// start from 0
	sum := make([][]int, len(matrix))
	for i := range sum {
		sum[i] = make([]int, len(matrix[i]))
	}
	// start from point (i, j), extend from left to right
	// then from up to down
	for x := 0; x < len(matrix); x++ {
		for y := 0; y < len(matrix[x]); y++ {
			if matrix[x][y] == k {
				return k
			}
			//fmt.Println("Start from:(", x, ",", y, ")")
			for i := x; i < len(matrix); i++ {
				for j := y; j < len(matrix[x]); j++ {
					curSum := matrix[i][j]
					if j != y {
						curSum += sum[i][j-1]
					}
					if i != x {
						curSum += sum[i-1][j]
					}
					if j != y && i != x {
						curSum -= sum[i-1][j-1]
					}
					//fmt.Println("---- sum(", i, ",", j, ")=", curSum, "sum:", sum)
					if curSum == k {

						return k
					}
					if curSum > maxValue && curSum < k {
						maxValue = curSum
					}
					sum[i][j] = curSum
				}
			}

		}
	}
	return maxValue
}

func Test_maxSumSubmatrix(t *testing.T) {
	type args struct {
		matrix [][]int
		k      int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				matrix: [][]int{{2, 2, -1}},
				k:      3,
			},
			want: 3,
		},
		{
			args: args{
				matrix: [][]int{{1, 0, 1}, {0, -2, 3}},
				k:      2,
			},
			want: 2,
		},
		{
			args: args{
				matrix: [][]int{{1, 0, 0}, {0, -2, 3}},
				k:      4,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxSumSubmatrix(tt.args.matrix, tt.args.k); got != tt.want {
				t.Errorf("maxSumSubmatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}
