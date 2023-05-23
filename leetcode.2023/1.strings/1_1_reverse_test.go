package strings_1_test

import "testing"

func reverseRune(r []rune, begin int, end int) {
	for begin < end {
		t := r[begin]
		r[begin] = r[end]
		r[end] = t
		begin++
		end--
	}
}

func moveString(s string, offset int) string {
	if offset < 0 || offset >= len(s) {
		panic("invalid offset")
	}
	r := []rune(s)
	reverseRune(r, 0, offset)
	reverseRune(r, offset+1, len(r)-1)
	reverseRune(r, 0, len(r)-1)
	return string(r)
}

func reverseWords(s string) string {
	r := []rune(s)
	startOffset := 0
	for i, c := range r {
		if c == ' ' {
			reverseRune(r, startOffset, i)
			startOffset = i + 1
		}
	}
	reverseRune(r, 0, len(r)-1)
	return string(r)
}

func Test_MoveString(t *testing.T) {
	testCases := []struct {
		s      string
		offset int
		expect string
	}{
		{"s", 0, "s"},
		{"abcdefg", 2, "defgabc"},
	}

	for _, tc := range testCases {
		got := moveString(tc.s, tc.offset)
		if got != tc.expect {
			t.Fatalf("Input: %s Expect: %s Got: %s", tc.s, tc.expect, got)
		}
	}
}

func Test_ReverseWords(t *testing.T) {
	got := reverseWords("I am a student.")
	if got != "student. a am I" {
		t.Fatalf("Bad output: %s", got)
	}
}
