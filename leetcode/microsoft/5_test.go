package microsoft

import (
	"testing"
)

// 采用编程之法1.6 最长回文子串的中心扩展法解决。另外也可以通过manacher解决，不掌握
func longestPalindrome(s string) string {
	if len(s) == 0 {
		return s
	}
	var maxStart int
	var maxLen int
	var j, length int
	for i := range s {
		// odd
		for j = 0; (i-j >= 0) && (i+j < len(s)); j++ {
			if s[i-j] != s[i+j] {
				break
			}
			length = j*2 + 1
		}
		if length > maxLen {
			maxStart = i - length/2
			maxLen = length
		}
		// even
		for j = 0; (i-j >= 0) && (i+j+1 < len(s)); j++ {
			if s[i-j] != s[i+j+1] {
				break
			}
			length = j*2 + 2
		}
		if length > maxLen {
			maxStart = i - length/2 + 1
			maxLen = length
		}
	}
	return s[maxStart : maxStart+maxLen]
}

// Bad manacher
func longestPalindromeBasicManacher(s string) string {
	if len(s) == 0 {
		return s
	}
	m := make([]byte, len(s)*2+1)
	m[0] = '#'
	for i := range s {
		m[i*2+1] = s[i]
		m[i*2+2] = '#'
	}
	var maxStart int
	var maxLen int
	var length int
	var j int
	for i := range m {
		for j = 0; (i-j >= 0) && (i+j < len(m)); j++ {
			if m[i-j] != m[i+j] {
				break
			}
			length = j*2 + 1
		}
		if length > maxLen {
			maxLen = length
			maxStart = i - length/2
		}
	}
	buf := make([]byte, 0, maxLen)
	for _, c := range m[maxStart : maxStart+maxLen] {
		if c != '#' {
			buf = append(buf, c)
		}
	}
	return string(buf)
}

func Test_longestPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "babad",
			want: "bab",
		},
		{
			in:   "abac",
			want: "aba",
		},
		{
			in:   "cbbd",
			want: "bb",
		},
		{
			in:   "ac",
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := longestPalindromeBasicManacher(tt.in); got != tt.want {
				t.Errorf("longestPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
