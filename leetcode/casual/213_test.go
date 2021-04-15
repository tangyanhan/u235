package casual

import (
	"testing"
)

func robHelper(nums, dp []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}

	dp[0] = nums[0]
	if nums[1] >= nums[0] {
		dp[1] = nums[1]
	}
	for i := 2; i < len(nums); i++ {
		robThis := dp[i-2] + nums[i]
		if robThis > dp[i-1] {
			dp[i] = robThis
		} else {
			dp[i] = dp[i-1]
		}
	}
	return dp[len(dp)-1]
}

func rob(nums []int) int {
	if len(nums) < 2 {
		return robHelper(nums, nil)
	}
	dp := make([]int, len(nums)-1)
	robFirst := robHelper(nums[:len(nums)-1], dp)
	robLast := robHelper(nums[1:], dp)
	if robFirst > robLast {
		return robFirst
	}
	return robLast
}

func Test_rob(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			nums: []int{1, 3, 1},
			want: 3,
		},
		{
			nums: []int{1, 2, 3, 1},
			want: 4,
		},
		{
			nums: []int{1},
			want: 1,
		},
		{
			nums: []int{1, 3, 1, 3, 100},
			want: 103,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rob(tt.nums); got != tt.want {
				t.Errorf("rob() = %v, want %v", got, tt.want)
			}
		})
	}
}
