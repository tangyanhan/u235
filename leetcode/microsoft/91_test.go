package microsoft

import (
	"testing"
)

func numDecodings(s string) int {
	if len(s) == 0 || s[0] == '0' {
		return 0
	}
	prev, curr := 1, 1
	for i := 1; i < len(s); i++ {
		tmp := curr
		if s[i] == '0' {
			if s[i-1] == '1' || s[i-1] == '2' {
				curr = prev
				continue
			}
			return 0
		}
		if s[i-1] == '1' || (s[i-1] == '2' && s[i] > '0' && s[i] < '7') {
			curr = curr + prev
		}
		prev = tmp
	}
	return curr
}

func Test_numDecodings(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		s    string
		want int
	}{
		{
			s:    "0",
			want: 0,
		},
		{
			s:    "99",
			want: 1,
		},
		{
			s:    "226",
			want: 3,
		},
		{
			s:    "2262",
			want: 3,
		},
		{
			s:    "22626",
			want: 6,
		},
		{
			s:    "10",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numDecodings(tt.s); got != tt.want {
				t.Errorf("numDecodings() = %v, want %v", got, tt.want)
			}
		})
	}
}
