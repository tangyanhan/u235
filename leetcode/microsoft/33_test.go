package microsoft

import "testing"

func binarySearch(nums []int, target int) int {
	low, high := 0, len(nums)
	for mid := (low + high) / 2; low < high; mid = (low + high) / 2 {
		switch {
		case nums[mid] == target:
			return mid
		case nums[mid] > target:
			high = mid
		default:
			low = mid + 1
		}
	}
	return -1
}

func search(nums []int, target int) int {
	var divIdx int
	for divIdx = 1; divIdx < len(nums) && nums[divIdx-1] < nums[divIdx]; divIdx++ {
	}
	// 数组存在旋转
	if divIdx < len(nums) {
		if target < nums[0] {
			if v := binarySearch(nums[divIdx:], target); v != -1 {
				return v + divIdx
			}
			return -1
		}
		return binarySearch(nums[0:divIdx], target)
	}

	// 数组太短或不存在旋转
	return binarySearch(nums, target)
}

func Test_search(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				nums:   []int{4, 5, 6, 7, 0, 1, 2},
				target: 0,
			},
			want: 4,
		},
		{
			args: args{
				nums:   []int{4, 5, 6, 7, 0, 1, 2},
				target: 5,
			},
			want: 1,
		},
		{
			args: args{
				nums:   []int{1},
				target: 0,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := search(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("search() = %v, want %v", got, tt.want)
			}
		})
	}
}
