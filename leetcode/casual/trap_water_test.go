package casual

import (
	"testing"
)

func trap(height []int) int {
	if len(height) == 0 {
		return 0
	}
	left := make([]int, len(height))
	left[0] = height[0]
	for i := 1; i < len(left); i++ {
		if height[i] > left[i-1] {
			left[i] = height[i]
		} else {
			left[i] = left[i-1]
		}
	}
	right := make([]int, len(height))
	right[len(right)-1] = height[len(right)-1]
	for i := len(right) - 2; i >= 0; i-- {
		if height[i] > right[i+1] {
			right[i] = height[i]
		} else {
			right[i] = right[i+1]
		}
	}
	var total int
	for i := 0; i < len(right); i++ {
		minHeight := left[i]
		if minHeight > right[i] {
			minHeight = right[i]
		}
		if height[i] < minHeight {
			total += minHeight - height[i]
		}
	}
	return total
}

func Test_trap(t *testing.T) {
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{
			height: []int{1, 0, 1},
			want:   1,
		},
		{
			height: []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1},
			want:   6,
		},
		{
			height: []int{0, 0, 0, 0},
			want:   0,
		},
		{
			height: []int{0, 0, 0, 1},
			want:   0,
		},
		{
			height: []int{5, 4, 1, 2},
			want:   1,
		},
		{
			height: []int{5, 2, 1, 2, 1, 5},
			want:   14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trap(tt.height); got != tt.want {
				t.Errorf("trap() = %v, want %v", got, tt.want)
			}
		})
	}
}
