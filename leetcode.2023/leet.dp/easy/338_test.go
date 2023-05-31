package easy_test

import (
	"log"
	"reflect"
	"testing"
)

func count1(n int) int {
	const tmpl = 0x01
	var count int
	for ; n != 0; n >>= 1 {
		if (tmpl & n) != 0 {
			count++
		}
	}

	return count
}

func countBits(n int) []int {
	ret := make([]int, n+1)

	for i := 0; i <= n; i++ {
		ret[i] = count1(i)
	}

	log.Println(ret)
	return ret
}

func Test_countBits(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{"2", 2, []int{0, 1, 1}},
		{"5", 5, []int{0, 1, 1, 2, 1, 2}},
		{"6", 5, []int{0, 1, 1, 2, 1, 2, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countBits(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
