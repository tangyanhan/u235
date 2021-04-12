package casual

import (
	"bytes"
	"sort"
	"strconv"
	"testing"
)

func largestNumber(nums []int) string {
	base := make(map[int]int)
	getBase := func(v int) int {
		b, ok := base[v]
		if ok {
			return b
		}
		b = 10
		for v >= b {
			b *= 10
		}
		base[v] = b
		return b
	}
	sort.Slice(nums, func(i, j int) bool {
		bI, bJ := getBase(nums[i]), getBase(nums[j])
		comboI := nums[i]*bJ + nums[j]
		comboJ := nums[j]*bI + nums[i]
		return comboI > comboJ
	})
	if nums[0] == 0 {
		return "0"
	}
	var buf bytes.Buffer
	for _, v := range nums {
		buf.WriteString(strconv.Itoa(v))
	}
	return buf.String()
}

func Test_largestNumber(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want string
	}{
		{
			nums: []int{10, 2},
			want: "210",
		},
		{
			nums: []int{3, 30, 34, 5, 9},
			want: "9534330",
		},
		{
			nums: []int{0, 0, 0, 0},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := largestNumber(tt.nums); got != tt.want {
				t.Errorf("largestNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
