package set_bilibili

import (
	"testing"
)

// https://leetcode-cn.com/problems/longest-substring-without-repeating-characters/
func lengthOfLongestSubstring(s string) int {
	var maxLen int
	var lastOffset int
	offsets := make([]int, 128)
	for i, r := range s {
		idx := int(r)
		repeatAt := offsets[idx] - 1
		if repeatAt >= lastOffset {
			l := i - lastOffset
			if l > maxLen {
				maxLen = l
			}
			lastOffset = repeatAt + 1
		}
		offsets[idx] = i + 1
	}
	if len(s) != lastOffset {
		l := len(s) - lastOffset
		if l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

func Test_lengthOfLongestSubstring(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{"abcabcbb"},
			want: 3,
		},
		{
			args: args{"bbbbbbb"},
			want: 1,
		},
		{
			args: args{"pwwkew"},
			want: 3,
		},
		{
			args: args{"cdd"},
			want: 2,
		},
		{
			args: args{" "},
			want: 1,
		},
		{
			args: args{""},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lengthOfLongestSubstring(tt.args.s); got != tt.want {
				t.Errorf("lengthOfLongestSubstring() input=%s = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}
