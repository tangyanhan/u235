package microsoft

import (
	"testing"
)

func compareVersion(va string, vb string) int {
	next := func(s string, i int) (int, int) {
		if i == -1 || i >= len(s) {
			return -1, 0
		}
		for ; i < len(s) && s[i] == '.'; i++ {
		}
		for ; i < len(s) && s[i] == '0'; i++ {
		}
		start := i
		for ; i < len(s) && s[i] != '.'; i++ {

		}
		num := 0
		factor := 1
		for j := i - 1; j >= start; j-- {
			num += factor * int(s[j]-'0')
			factor *= 10
		}

		if i < len(s) || i != start {
			return i, num
		}
		return -1, 0
	}

	var ia, ib int
	var a, b int
	for ia != -1 || ib != -1 {
		ia, a = next(va, ia)
		ib, b = next(vb, ib)
		if a != b {
			if a > b {
				return 1
			}
			return -1
		}
	}
	return 0
}

func Test_compareVersion(t *testing.T) {
	tests := []struct {
		name string
		ver  []string
		want int
	}{
		{
			ver:  []string{"0.0.1", "0.0.1.0"},
			want: 0,
		},
		{
			ver:  []string{"2.3.1", "2.3"},
			want: 1,
		},
		{
			ver:  []string{"2.3.1", "2.3.2"},
			want: -1,
		},
		{
			ver:  []string{"0.1", "1.0"},
			want: -1,
		},
		{
			ver:  []string{"1.0", "1.0.0"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareVersion(tt.ver[0], tt.ver[1]); got != tt.want {
				t.Errorf("compareVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
