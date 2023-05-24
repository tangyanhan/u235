package easy_test

import "testing"

func isValidChar(c byte) bool {
	return (c >= byte('a') && c <= byte('z')) || (c >= byte('A') && c <= byte('Z')) || (c >= byte('0') && c <= byte('9'))
}

func isEqualIgnoreCase(c1, c2 byte) bool {
	x1 := c1 - 'A'
	if x1 >= 26 {
		x1 = c1 - 'a'
	}
	x2 := c2 - 'A'
	if x2 >= 26 {
		x2 = c2 - 'a'
	}
	return x1 == x2
}

func isPalindrome(s string) bool {
	head, tail := 0, len(s)-1
	for head < tail {
		for ; head < tail && !isValidChar(s[head]); head++ {
		}

		for ; head < tail && !isValidChar(s[tail]); tail-- {
		}

		if !isEqualIgnoreCase(s[head], s[tail]) {
			return false
		}
		head++
		tail--
	}
	return true
}

func Test_IsPalindrome(t *testing.T) {
	testCases := []struct {
		s      string
		expect bool
	}{
		{"A man, a plan, a canal: Panama", true},
		{"race a car", false},
		{"", true},
		{"0P", false},
	}

	for _, tc := range testCases {
		if got := isPalindrome(tc.s); got != tc.expect {
			t.Fatalf("%s Expect:%t Got:%t", tc.s, tc.expect, got)
		}
	}
}
