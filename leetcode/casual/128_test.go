package casual

import (
	"testing"
)

func longestConsecutive(nums []int) int {
	if len(nums) < 2 {
		return len(nums)
	}
	size := make([]int, len(nums))
	parent := make([]int, len(nums))
	idxMap := make(map[int]int)
	for i := range parent {
		idxMap[nums[i]] = i
		parent[i] = i
		size[i] = 1
	}
	find := func(p int) int {
		for p != parent[p] {
			parent[p] = parent[parent[p]]
			size[p] = size[parent[p]]
			p = parent[p]
		}
		return p
	}
	maxSize := 1
	union := func(x, y int) {
		px, py := find(x), find(y)
		if px != py {
			parent[px] = py
			newSize := size[py] + size[px]
			if maxSize < newSize {
				maxSize = newSize
			}
			size[px] = newSize
			size[py] = newSize
		}
	}

	for v, i := range idxMap {
		low, high := v-1, v+1
		if iLow, ok := idxMap[low]; ok {
			union(iLow, i)
		}
		if iHigh, ok := idxMap[high]; ok {
			union(iHigh, i)
		}
	}
	return maxSize
}

func Test_longestConsecutive(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				nums: []int{3, 7, 2, 5, 8, 4, 6, 0, 1},
			},
			want: 9,
		},
		{
			args: args{
				nums: []int{0, 3, 7, 2, 5, 8, 4, 6, 0, 1},
			},
			want: 9,
		},
		{
			args: args{
				nums: []int{1, 3, 5},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestConsecutive(tt.args.nums); got != tt.want {
				t.Errorf("longestConsecutive() = %v, want %v", got, tt.want)
			}
		})
	}
}
