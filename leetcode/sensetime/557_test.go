package sensetime

import "testing"

func reverseWords(s string) string {
	buf := []byte(s)
	reverse := func(start, end int) {
		for start < end {
			buf[start], buf[end] = buf[end], buf[start]
			start++
			end--
		}
	}
	start := 0
	for i, c := range buf {
		if c == ' ' {
			reverse(start, i-1)
			start = i + 1
		}
	}
	reverse(start, len(buf)-1)
	return string(buf)
}

func Test_reverseWords(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "Let's take LeetCode contest",
			want: "s'teL ekat edoCteeL tsetnoc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := reverseWords(tt.in); got != tt.want {
				t.Errorf("reverseWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
