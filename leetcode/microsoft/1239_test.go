package microsoft

import (
	"fmt"
	"testing"
)

func TestWalkArray(t *testing.T) {
	var m [4][4]int
	num := 1
	for delta := 0; delta < len(m); delta++ {
		i, j := 0, delta
		for i < len(m) && j < len(m) {
			m[i][j] = num
			i++
			j++
		}
		num++
	}
	for _, row := range m {
		for _, v := range row {
			fmt.Printf("%x ", v)
		}
		fmt.Println("")
	}
}

func maxLength(arr []string) int {
	return 0
}

func Test_maxLength(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want int
	}{
		{
			in:   []string{"cha", "r", "act", "ers"},
			want: 6,
		},
		{
			in:   []string{"a", "b", "bc", "bd", "ef", "ac", "xyz"},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxLength(tt.in); got != tt.want {
				t.Errorf("maxLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
