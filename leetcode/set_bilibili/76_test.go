package set_bilibili

import "testing"

// https://leetcode-cn.com/problems/minimum-window-substring/
// 1. When a s' contains t, it must starts and ends with one character in t
// 2. s' may contains several characters not needed/repeated
// wip status
func minWindow(s, t string) string {
	tm := make(map[byte]bool)
	for i := 0; i < len(t); i++ {
		tm[t[i]] = true
	}
	var low, high int
	for i := 0; i < len(s); i++ {
		if tm[t[i]] {
			low = i
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if tm[s[i]] {
			high = i
			break
		}
	}
	if low >= high {
		return ""
	}

	return ""
}

func Test_minWindow(t *testing.T) {
	type args struct {
		s string
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{s: "ADOBECODEBANC", t: "ABC"},
			want: "BANC",
		},
		{
			args: args{s: "a", t: "a"},
			want: "a",
		},
		{
			args: args{s: "xyz", t: "a"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minWindow(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("minWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
