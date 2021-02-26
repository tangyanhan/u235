package sensetime

import (
	"reflect"
	"testing"
)

func singleNumber2(nums []int) []int {
	var x int
	for _, v := range nums {
		x ^= v
	}
	// now x= a^b
	var mask int
	for i := 1; i < 32; i++ {
		mask = 0x01 << i
		if mask&x != 0 {
			break
		}
	}
	var xa, xb int
	for _, v := range nums {
		if mask&v != 0 {
			xa ^= v
		} else {
			xb ^= v
		}
	}
	return []int{xa, xb}
}

func Test_singleNumber2(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{nums: []int{1, 2, 1, 3, 2, 5}},
			want: []int{3, 5},
		},
		{
			args: args{nums: []int{1193730082, 587035181, -814709193, 1676831308, -511259610, 284593787, -2058511940, 1970250135, -814709193, -1435587299, 1308886332, -1435587299, 1676831308, 1403943960, -421534159, -528369977, -2058511940, 1636287980, -1874234027, 197290672, 1976318504, -511259610, 1308886332, 336663447, 1636287980, 197290672, 1970250135, 1976318504, 959128864, 284593787, -528369977, -1874234027, 587035181, -421534159, -786223891, 933046536, 959112204, 336663447, 933046536, 959112204, 1193730082, -786223891}},
			want: []int{1403943960, 959128864},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := singleNumber2(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("singleNumber2() = %v, want %v", got, tt.want)
			}
		})
	}
}
