package conquor

func simple(nums []int) int {
	cnt := make(map[int]int)
	maxCount, maxNum := 0, 0
	for _, v := range nums {
		cnt[v]++
		if cnt[v] > maxCount {
			maxCount = cnt[v]
			maxNum = v
		}
	}
	return maxNum
}

func majorityElement(nums []int) int {
	count, result := 0, 0
	for _, v := range nums {
		if count == 0 {
			count = 1
			result = v
		} else if v == result {
			count++
		} else {
			count--
		}
	}
	return result
}
