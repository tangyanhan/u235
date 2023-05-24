package easy_test

import "testing"

var psIdx int
var parenthStack [5001]byte

func pushPS(r rune) {
	parenthStack[psIdx] = byte(r)
	psIdx++
}

func isValid(s string) bool {
	return false
}

func Test_ParenthesesIsValid(t *testing.T) {
	testCases := []struct {
		s      string
		expect bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"(]", false},
	}
	for _, tc := range testCases {
		got := isValid(tc.s)
		if got != tc.expect {
			t.Fatalf("%s Expect:%t Got:%t", tc.s, tc.expect, got)
		}
	}
}
