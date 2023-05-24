package easy_test

import (
	"strings"
	"testing"
)

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	for i, r := range strs[0] {
		b := byte(r)
		for _, s := range strs[1:] {
			if i >= len(s) || s[i] != b {
				return strs[0][:i]
			}
		}
	}

	return strs[0]
}

func TestLCP(t *testing.T) {
	testCases := []struct {
		strs   []string
		expect string
	}{
		{[]string{"flower", "flow", "flight"}, "fl"},
		{[]string{"dog", "racecar", "car"}, ""},
		{[]string{"", "abc", ""}, ""},
	}

	for _, tc := range testCases {
		got := longestCommonPrefix(tc.strs)
		if got != tc.expect {
			t.Fatalf("Input=%s Expect:%s Got:%s", strings.Join(tc.strs, ","), tc.expect, got)
		}
	}
}
