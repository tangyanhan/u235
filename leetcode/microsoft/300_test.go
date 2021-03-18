package microsoft

import "testing"

func lengthOfLIS(nums []int) int {
	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
	}
	var ans int
	for i := 1; i < len(dp); i++ {
		for j := i - 1; j >= 0; j-- {
			if nums[i] > nums[j] {
				if dp[i] < dp[j]+1 {
					dp[i] = dp[j] + 1
				}
			}
		}
	}
	for _, v := range dp {
		if ans < v {
			ans = v
		}
	}
	return ans
}

func lengthOfLISPatience(nums []int) int {
	var piles [2500]int
	var pileIdx int
	searchPile := func(v int) {
		low, high := 0, pileIdx
		for mid := (low + high) / 2; low < high; mid = (low + high) / 2 {
			if piles[mid] >= v {
				high = mid
			} else if piles[mid] < v {
				low = mid + 1
			}
		}
		if low == pileIdx {
			piles[pileIdx] = v
			pileIdx++
		} else {
			piles[low] = v
		}
	}
	for _, v := range nums {
		searchPile(v)
	}
	return pileIdx
}

func Test_lengthOfLIS(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{
			in:   []int{0, 1, 0, 3, 2, 3},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lengthOfLISPatience(tt.in); got != tt.want {
				t.Errorf("lengthOfLIS() = %v, want %v", got, tt.want)
			}
		})
	}
}
