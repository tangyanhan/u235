package microsoft

import (
	"testing"
)

func checkSubarraySumSlow(nums []int, k int) bool {
	for i := 1; i < len(nums); i++ {
		sum := nums[i]
		for j := i - 1; j >= 0; j-- {
			sum += nums[j]
			if (k == 0 && sum == 0) || (k != 0 && sum%k == 0) {
				return true
			}
		}
	}
	return false
}

func Test_checkSubarraySum(t *testing.T) {
	type args struct {
		nums []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				nums: []int{23, 2, 6, 4, 7},
				k:    6,
			},
			want: true,
		},
		{
			args: args{
				nums: []int{1, 0, 1, 0, 1},
				k:    4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkSubarraySumSlow(tt.args.nums, tt.args.k); got != tt.want {
				t.Errorf("checkSubarraySum() = %v, want %v", got, tt.want)
			}
		})
	}
}
