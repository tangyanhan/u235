package microsoft

import (
	"testing"
)

func maxArea(height []int) int {
	// 双指针法，从左右两侧开始寻找。已知目前收益为S，任一挪动left,right，
	// 都必须让更短的height[i]变大才有可能弥补宽度的减少，两者相等时任意移动一个即可
	left, right := 0, len(height)-1
	for left < right && height[left] == 0 {
		left++
	}
	for right > left && height[right] == 0 {
		right--
	}

	var maxArea int
	for left < right {
		minHeight := height[left]
		width := right - left
		if minHeight > height[right] {
			minHeight = height[right]
			right--
		} else {
			left++
		}

		prod := minHeight * width
		if prod > maxArea {
			maxArea = prod
		}
	}

	return maxArea
}

func Test_maxArea(t *testing.T) {
	type args struct {
		height []int
	}
	tests := []struct {
		name   string
		height []int
		want   int
	}{
		{
			height: []int{1, 8, 6, 2, 5, 4, 8, 3, 7},
			want:   49,
		},
		{
			height: []int{4, 3, 2, 1, 4},
			want:   16,
		},
		{
			height: []int{1, 2, 1},
			want:   2,
		},
		{
			height: []int{2, 3, 4, 5, 18, 17, 6},
			want:   17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxArea(tt.height); got != tt.want {
				t.Errorf("maxArea() = %v, want %v", got, tt.want)
			}
		})
	}
}
