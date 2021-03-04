package set_bilibili

import "testing"

func titleToNumber(s string) int {
	factor := 1
	var num int
	for i := len(s) - 1; i >= 0; i-- {
		v := int(s[i]-'A') + 1
		num += v * factor
		factor *= 26
	}
	return num
}

func Test_titleToNumber(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		in   string
		want int
	}{
		{
			in:   "A",
			want: 1,
		},
		{
			in:   "AB",
			want: 28,
		},
		{
			in:   "ZY",
			want: 701,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := titleToNumber(tt.in); got != tt.want {
				t.Errorf("titleToNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
