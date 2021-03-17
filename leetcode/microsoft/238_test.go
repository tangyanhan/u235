package microsoft

import (
	"reflect"
	"testing"
)

func productExceptSelf(nums []int) []int {
	left := make([]int, len(nums)+1)
	left[0] = 1
	for i := 0; i < len(nums); i++ {
		left[i+1] = left[i] * nums[i]
	}
	right := make([]int, len(nums)+1)
	right[len(right)-1] = 1

	for i := len(right) - 2; i >= 0; i-- {
		right[i] = right[i+1] * nums[i]
	}
	for i := 0; i < len(nums); i++ {
		nums[i] = left[i] * right[i+1]
	}
	return nums
}

func Test_productExceptSelf(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{
			nums: []int{-1, 1, 0, -3, 3},
			want: []int{0, 0, 9, 0, 0},
		},
		{
			nums: []int{1, 2, 3, 4},
			want: []int{24, 12, 8, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := productExceptSelf(tt.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productExceptSelf() = %v, want %v", got, tt.want)
			}
		})
	}
}
