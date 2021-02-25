package dp

func moveZeroes(nums []int) {
	var nonZero int
	for i, v := range nums {
		if v != 0 {
			if nonZero != i {
				nums[nonZero] = v
			}
			nonZero++
		}
	}
	for i := nonZero; i < len(nums); i++ {
		nums[i] = 0
	}
}
