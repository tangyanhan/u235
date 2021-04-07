package casual

import "testing"

func search(nums []int, target int) bool {
	low, high := 0, len(nums)
	for mid := (low + high) / 2; low < high; mid = (low + high) / 2 {
		if nums[mid] == target {
			return true
		}
		// 判断单调递增不能使用<=，因为可能出现 1, 2, 1, 2, 1导致判断失误
		isIncreasing := nums[low] < nums[mid] && nums[mid] < nums[high-1]
		if isIncreasing {
			if nums[mid] > target {
				high = mid
			} else {
				low = mid + 1
			}
		} else {
			return search(nums[low:mid], target) || search(nums[mid+1:high], target)
		}
	}
	return false
}

func Test_search(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				nums:   []int{5, 6, 7, 1, 2, 3, 4},
				target: 6,
			},
			want: true,
		},
		{
			args: args{
				nums:   []int{5, 6, 7, 1, 2, 3, 4},
				target: 3,
			},
			want: true,
		},
		{
			args: args{
				nums:   []int{5, 6, 7, 1, 2, 3, 4},
				target: 9,
			},
			want: false,
		},
		{
			args: args{
				nums:   []int{1},
				target: 1,
			},
			want: true,
		},
		{
			args: args{
				nums:   []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1},
				target: 2,
			},
			want: true,
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
