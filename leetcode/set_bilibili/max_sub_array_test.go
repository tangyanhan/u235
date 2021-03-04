package set_bilibili

func maxSubArray(nums []int) int {
	ans := nums[0]
	for i := 1; i < len(nums); i++ {
		nums[i] = max(nums[i], nums[i]+nums[i-1])
		if ans < nums[i] {
			ans = nums[i]
		}
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
