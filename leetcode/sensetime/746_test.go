package sensetime

import (
	"testing"
)

func minCostClimbingStairs(cost []int) int {
	data := make([]int, len(cost))
	for i, v := range cost {
		data[i] = v
	}
	min := func(i, j int) int {
		if i < 0 {
			return 0
		}
		if data[i] < data[j] {
			return data[i]
		}
		return data[j]
	}
	for i := 1; i < len(cost); i++ {
		data[i] = cost[i] + min(i-2, i-1)
	}
	return min(len(cost)-2, len(cost)-1)
}

func Test_minCostClimbingStairs(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		cost []int
		want int
	}{
		{
			cost: []int{10, 15, 20},
			want: 15,
		},
		{
			cost: []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minCostClimbingStairs(tt.cost); got != tt.want {
				t.Errorf("minCostClimbingStairs() = %v, want %v", got, tt.want)
			}
		})
	}
}
