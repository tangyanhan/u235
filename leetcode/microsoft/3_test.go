package microsoft

func lengthOfLongestSubstring(s string) int {
	var maxLen int
	var lastOffset int
	offsets := make([]int, 128)
	for i, r := range s {
		idx := int(r)
		repeatAt := offsets[idx] - 1
		if repeatAt >= lastOffset {
			l := i - lastOffset
			if l > maxLen {
				maxLen = l
			}
			lastOffset = repeatAt + 1
		}
		offsets[idx] = i + 1
	}
	if len(s) != lastOffset {
		l := len(s) - lastOffset
		if l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}
