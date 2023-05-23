package strings_1_test

import "testing"

func StringCharacterContains(a string, b string) bool {
	var sigA uint64
	for _, c := range a {
		offset := c - 'A'
		sigA |= 1 << offset
	}

	var sigB uint64
	for _, c := range b {
		offset := c - 'A'
		bitB := uint64(1 << offset)
		if bitB&sigA == 0 {
			return false
		}
		sigB |= bitB
	}
	return sigA == sigB
}

func Test_StringCharacterContains(t *testing.T) {
	testCases := []struct {
		a      string
		b      string
		expect bool
	}{
		{"ABD", "BAD", true},
		{"ABD", "BA", false},
		{"ABD", "BADE", false},
	}

	for _, tc := range testCases {
		got := StringCharacterContains(tc.a, tc.b)
		if got != tc.expect {
			t.Fatalf("%s ## %s Expect: %t Got: %t", tc.a, tc.b, tc.expect, got)
		}
	}
}
