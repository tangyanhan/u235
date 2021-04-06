package casual

import (
	"testing"
)

func removeDuplicates(nums []int) int {
	var insertIdx int
	insertNum := func(v int) {
		nums[insertIdx] = v
		insertIdx++
	}
	for i := 0; i < len(nums); {
		curr := nums[i]
		j := i + 1
		for ; j < len(nums) && nums[j] == curr; j++ {

		}
		length := j - i
		for k := 0; k < length && k < 2; k++ {
			insertNum(curr)
		}
		i = j
	}
	return insertIdx
}

func Test_removeDuplicates(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			nums: []int{0, 1, 2, 3, 4, 5},
			want: 6,
		},
		{
			nums: []int{0, 0, 1, 1, 1, 2, 3, 4},
			want: 7,
		},
		{
			nums: []int{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicates(tt.nums); got != tt.want {
				t.Errorf("removeDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
