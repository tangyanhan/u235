package microsoft

func maxProduct(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	dp := make([][2]int, len(nums))
	dp[0][0] = nums[0]
	dp[0][1] = nums[0]
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	maxProd := nums[0]
	for i := 1; i < len(dp); i++ {
		prodMax := dp[i-1][0] * nums[i]
		prodMin := dp[i-1][1] * nums[i]
		dp[i][0] = max(prodMax, max(prodMin, nums[i]))
		dp[i][1] = min(prodMin, min(prodMax, nums[i]))
		if dp[i][0] > maxProd {
			maxProd = dp[i][0]
		}
	}
	return maxProd
}
