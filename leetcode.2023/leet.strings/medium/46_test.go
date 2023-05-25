package medium_test

import (
	"bytes"
	"sort"
	"strconv"
	"testing"
)

func next(nums []int) bool {
	var i int
	for i = len(nums) - 2; i >= 0 && nums[i] >= nums[i+1]; i-- {
	}

	if i < 0 {
		return false
	}

	var k int
	for k = len(nums) - 1; k > i && nums[k] <= nums[i]; k-- {
	}

	nums[i], nums[k] = nums[k], nums[i]
	i++
	k = len(nums) - 1
	for i < k {
		nums[i], nums[k] = nums[k], nums[i]
		i++
		k--
	}
	return true
}

func permute(nums []int) [][]int {
	sort.Ints(nums)
	result := make([][]int, 0, 999999)
	item := make([]int, len(nums))
	copy(item, nums)
	result = append(result, item)
	for next(nums) {
		item := make([]int, len(nums))
		copy(item, nums)
		result = append(result, item)
	}
	return result
}

func JoinSliceString(s []int, separator string) string {
	var buf bytes.Buffer
	for i, n := range s {
		buf.WriteString(strconv.FormatInt(int64(n), 10))
		if i != len(s)-1 {
			buf.WriteString(separator)
		}
	}
	return buf.String()
}

func Test_Permute(t *testing.T) {
	testCases := []struct {
		in     []int
		expect [][]int
	}{
		{[]int{0, 1}, [][]int{{0, 1}, {1, 0}}},
		{[]int{0, -1, 1}, [][]int{{0, 1}, {1, 0}}},
	}

	for _, tc := range testCases {
		got := permute(tc.in)
		t.Log(len(got))
		for i, nums := range got {
			t.Logf("#%d %s", i, JoinSliceString(nums, ","))
		}
	}
}
