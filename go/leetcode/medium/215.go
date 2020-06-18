package medium

import (
	"container/heap"
)

type MinSlice []int

func (m *MinSlice) Len() int {
	return len(*m)
}

func (m *MinSlice) Less(i, j int) bool {
	return (*m)[i] < (*m)[j]
}

func (m *MinSlice) Swap(i, j int) {
	s := *m
	s[i], s[j] = s[j], s[i]
}

func (m *MinSlice) Push(v interface{}) {
}

func (m *MinSlice) Pop() interface{} {
	return 0
}

func findKthLargest(nums []int, k int) int {
	h := MinSlice(nums[:k])
	heap.Init(&h)
	for _, v := range nums[k:] {
		if h[0] < v {
			h[0] = v
			heap.Fix(&h, 0)
		}
	}
	return h[0]
}
