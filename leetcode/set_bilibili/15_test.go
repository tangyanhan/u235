package set_bilibili

import (
	"reflect"
	"sort"
	"testing"
)

func binarySearch(s []int, target int) int {
	if len(s) == 0 {
		return -1
	}
	if s[0] > target || s[len(s)-1] < target {
		return -1
	}
	low, high := 0, len(s)

	for mid := (low + high) / 2; low < high; mid = (low + high) / 2 {
		if s[mid] == target {
			return mid
		} else if s[mid] < target {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return -1
}

func TestBinSearch(t *testing.T) {
	t.Log(binarySearch([]int{1, 3}, 2))
	t.Log(binarySearch([]int{1, 3}, 3))
	t.Log(binarySearch([]int{1}, 1))
}

// https://leetcode-cn.com/problems/3sum/
func threeSum(nums []int) [][]int {
	sort.Ints(nums)
	m := make(map[int][2]int)
	for i := 0; i < len(nums); {
		v := nums[i]
		tuple := [2]int{i, 0}
		for i++; i < len(nums) && nums[i] == v; i++ {
		}
		tuple[1] = i
		m[v] = tuple
	}
	// Given any a+b, c = 0 - a - b, so if we simply take all a+b, O(N^2)
	lastA, lastB := 999999, 999999
	result := make([][]int, 0)
	for i, a := range nums {
		sum2 := 0 - a
		subset2 := nums[i+1:]
		for j, b := range subset2 {
			// Skip duplicate search
			if a == lastA && b == lastB {
				continue
			}
			lastA, lastB = a, b
			target := sum2 - b
			minIdx := i + j + 2
			if tuple, ok := m[target]; ok && minIdx < tuple[1] {
				if len(result) != 0 {
					lastVex := result[len(result)-1]
					if lastVex[0] >= a && lastVex[1] >= b {
						continue
					}
				}
				vex := []int{a, b, target}
				result = append(result, vex)
			}
		}
	}
	return result
}

func threeSumDbl(nums []int) [][]int {
	sort.Ints(nums)

	results := make([][]int, 0, len(nums))
	var i, low, high int
	nextLow := func() {
		v := nums[low]
		low++
		for low < high && nums[low] == v {
			low++
		}
	}
	nextHigh := func() {
		v := nums[high]
		high--
		for low < high && nums[high] == v {
			high--
		}
	}
	nextI := func() {
		v := nums[i]
		i++
		for i < len(nums) && nums[i] == v {
			i++
		}
	}
	for i < len(nums)-2 {
		low, high = i+1, len(nums)-1
		a := nums[i]
		remain := 0 - a
		if remain < a {
			nextI()
			continue
		}
		for low < high {
			sum := nums[low] + nums[high]
			switch {
			case sum == remain:
				results = append(results, []int{a, nums[low], nums[high]})
				nextLow()
				nextHigh()
			case sum < remain:
				nextLow()
			default:
				nextHigh()
			}
		}
		nextI()
	}
	return results
}

// https://leetcode-cn.com/problems/3sum/
func threeSumBin(nums []int) [][]int {
	sort.Ints(nums)
	// Given any a+b, c = 0 - a - b, so if we simply take all a+b, O(N^2)
	lastA, lastB := 999999, 999999
	result := make([][]int, 0)
	for i, a := range nums {
		sum2 := 0 - a
		subset2 := nums[i+1:]
		for j, b := range subset2 {
			// Skip duplicate search
			if a == lastA && b == lastB {
				continue
			}
			lastA, lastB = a, b
			target := sum2 - b
			subset1 := subset2[j+1:]
			if k := binarySearch(subset1, target); k != -1 {
				vex := []int{a, b, target}
				if len(result) != 0 {
					lastVex := result[len(result)-1]
					if lastVex[0] >= a && lastVex[1] >= b {
						continue
					}
				}
				result = append(result, vex)
			}
		}
	}
	return result
}

func Test_threeSum(t *testing.T) {
	type args struct {
		in0 []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			args: args{[]int{-1, 0, 1, 2, -1, -4}},
			want: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
		{
			args: args{[]int{-1, 0, 1}},
			want: [][]int{{-1, 0, 1}},
		},
		{
			args: args{[]int{}},
			want: [][]int{},
		},
		{
			args: args{[]int{0}},
			want: [][]int{},
		},
		{
			args: args{[]int{3, -2, 1, 0}},
			want: [][]int{},
		},
		{
			args: args{[]int{-4, -2, -2, -2, 0, 1, 2, 2, 2, 3, 3, 4, 4, 6, 6}},
			want: [][]int{{-4, -2, 6}, {-4, 0, 4}, {-4, 1, 3}, {-4, 2, 2}, {-2, -2, 4}, {-2, 0, 2}},
		},
		{
			args: args{[]int{2, -2, 0, 0, 0, 0, 0, 0, 0}},
			want: [][]int{{-2, 0, 2}, {0, 0, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := threeSumDbl(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threeSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
