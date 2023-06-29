package offer100_test

import (
	"testing"
)

func lengthOfLongestSubstring(s string) int {
	var offset [128]int
	for i := 0; i < len(offset); i++ {
		offset[i] = -1
	}

	bs := []byte(s)
	var start int
	var longestLength int
	for i, b := range bs {
		idx := int(b)
		//log.Println("i=", i, "start=", start, "idx=", idx, "previous offset=", offset[idx])
		// Hit a same offset within range
		if i > start && offset[idx] >= start {
			length := i - start
			if length > longestLength {
				longestLength = length
			}
			//log.Println("Hit same character at idx:", i, "Duplicate:", offset[idx], "Start=", start, "Length=", length, "String=", string(bs[start:i]))
			length = 1
			start = offset[idx] + 1
		}

		offset[idx] = i
	}

	lastLength := len(bs) - start
	if lastLength > longestLength {
		longestLength = lastLength
	}

	return longestLength
}

func Test_lengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"abcabcbb", 3},
		{"bbbbbb", 1},
		{"pwwkew", 3},
		{"abcdef", 6},
		{" ", 1},
		{"dvdf", 3},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := lengthOfLongestSubstring(tt.s); got != tt.want {
				t.Errorf("lengthOfLongestSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}
