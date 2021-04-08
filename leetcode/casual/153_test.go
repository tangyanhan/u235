package casual

import "testing"

func findMin(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	isAscending := nums[0] < nums[len(nums)-1]
	if isAscending {
		return nums[0]
	}
	mid := len(nums) / 2
	leftMin := findMin(nums[0:mid])
	rightMin := findMin(nums[mid:])
	if leftMin < rightMin {
		return leftMin
	}
	return rightMin
}

func Test_findMin(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			nums: []int{4, 5, 6, 7, 0, 1, 2},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMin(tt.nums); got != tt.want {
				t.Errorf("findMin() = %v, want %v", got, tt.want)
			}
		})
	}
}
