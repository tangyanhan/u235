package str

import (
	"bytes"
	"strings"
)

// ReverseString reverse string part
func ReverseString(s []rune, from, to int) {
	for from < to {
		t := s[from]
		s[from] = s[to]
		s[to] = t
		from++
		to--
	}
}

// LeftRotateString move string to left by offset m
func LeftRotateString(s string, m int) string {
	m %= len(s)
	r := []rune(s)
	ReverseString(r, 0, m-1)
	ReverseString(r, m, len(s)-1)
	ReverseString(r, 0, len(s)-1)
	return string(r)
}

// 题目：以单词为单位反转句子，句子单词以空格为间隔
func ReverseSentence(s string) string {
	p := strings.Split(s, " ")
	var buf bytes.Buffer
	for i := len(p) - 1; i >= 0; i-- {
		buf.WriteString(p[i])
		if i != 0 {
			buf.WriteRune(' ')
		}
	}
	return buf.String()
}
