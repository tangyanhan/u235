package microsoft

// 桶排序
func sortColors(nums []int) {
	var buckets [3]int
	for _, v := range nums {
		buckets[v]++
	}
	var i int
	for color, cnt := range buckets[:] {
		for j := 0; j < cnt; j++ {
			nums[i] = color
			i++
		}
	}
}
