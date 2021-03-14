package microsoft

import (
	"sort"
)

// twoSum with extra sort and space
func twoSum(nums []int, target int) []int {
	meta := make([]struct {
		value int
		pos   int
	}, len(nums))
	for i, v := range nums {
		meta[i].pos = i
		meta[i].value = v
	}
	sort.Slice(meta, func(i, j int) bool {
		return meta[i].value < meta[j].value
	})
	low, high := 0, len(nums)-1
	for low < high {
		delta := target - meta[low].value - meta[high].value
		if delta == 0 {
			return []int{meta[low].pos, meta[high].pos}
		}
		if delta < 0 {
			high--
		} else {
			low++
		}
	}
	return nil
}

// with extra map to reduce time to search
func twoSumMap(nums []int, target int) []int {
	m := make(map[int]int)
	half := target / 2
	for i, v := range nums {
		if v == half {
			pos, ok := m[v]
			if ok {
				return []int{pos, i}
			}
		}
		m[v] = i
	}
	for v, pos1 := range m {
		if pos2, ok := m[target-v]; ok && pos1 != pos2 {
			return []int{pos1, pos2}
		}
	}
	return nil
}
