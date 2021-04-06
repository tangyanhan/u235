package casual

import (
	"reflect"
	"testing"
)

func merge(nums1 []int, m int, nums2 []int, n int) {
	if n == 0 {
		return
	}
	if m == 0 {
		copy(nums1, nums2)
		return
	}
	var i int
	var j int
	var cnt int
	for j < n {
		if cnt >= m {
			copy(nums1[i:], nums2[j:n])
			return
		}
		if nums1[i] <= nums2[j] {
			i++
			cnt++
		} else {
			copy(nums1[i+1:], nums1[i:])
			nums1[i] = nums2[j]
			i++
			j++
		}
	}
}

func Test_merge(t *testing.T) {
	type args struct {
		nums1 []int
		m     int
		nums2 []int
		n     int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{
				nums1: []int{0, 2, 3, 0, 0},
				m:     3,
				nums2: []int{1, 4, 6},
				n:     2,
			},
			want: []int{0, 1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merge(tt.args.nums1, tt.args.m, tt.args.nums2, tt.args.n)
			if !reflect.DeepEqual(tt.args.nums1, tt.want) {
				t.Fatal(tt.args.nums1)
			}
		})
	}
}
