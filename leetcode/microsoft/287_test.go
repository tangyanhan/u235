package microsoft

import "testing"

// TODO: add the floyd circuit
func findDuplicate(nums []int) int {
	func findDuplicate(nums []int) int {
		var m [3*10000 + 3]int
		for _, v := range nums {
			m[v]++
			if m[v] != 1 {
				return v
			}
		}
		return -1
	}
}

func Test_findDuplicate(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{
			nums: []int{3, 1, 3, 4, 2},
			want: 3,
		},
		{
			nums: []int{1, 3, 4, 2, 2},
			want: 2,
		},
		{
			nums: []int{2,2,2,2,2},
			want: 2,
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findDuplicate(tt.nums); got != tt.want {
				t.Errorf("findDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
