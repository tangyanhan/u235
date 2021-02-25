package dp

import (
	"log"
	"testing"
)

func TestNonRepeatLen(t *testing.T) {
	t.Log(lengthOfLongestSubstring("abcabcbb"))
	t.Log(lengthOfLongestSubstring(" "))
}

type charMap struct {
	low  uint64
	high uint64
}

func (m *charMap) addChar(c byte) {
	if c > 63 {
		c -= 64
		m.high |= 0x1 << c
	} else {
		m.low |= 0x1 << c
	}
}

func (m *charMap) contains(c byte) bool {
	if c > 63 {
		c -= 64
		v := uint64(0x1) << c
		return m.high&v != 0
	}
	v := uint64(0x1) << c
	return m.low&v != 0
}

func (m *charMap) reset() {
	m.low, m.high = 0, 0
}

func lengthOfLongestSubstring(s string) int {
	var m charMap
	var maxLen, idx, i int
	bs := []byte(s)
	for i = 0; i < len(bs); i++ {
		c := bs[i]
		if m.contains(c) {
			curLen := i - idx
			if curLen > maxLen {
				maxLen = curLen
			}
			idx = i
			m.reset()
		}
		m.addChar(c)
	}
	log.Println("DBG - i=", i)
	lastLen := i - idx
	if lastLen > maxLen {
		maxLen = lastLen
	}
	return maxLen
}
