package strings_1_test

import (
	"log"
	"testing"
)

func reverseBytes(bs []byte, begin, end int) {
	for begin < end {
		bs[begin], bs[end] = bs[end], bs[begin]
		begin++
		end--
	}
}

func nextPermutation(bs []byte) bool {
	var i int
	for i = len(bs) - 2; i >= 0 && bs[i] >= bs[i+1]; i-- {
	}

	if i < 0 {
		return false
	}

	var k int
	for k = len(bs) - 1; k > i && bs[k] <= bs[i]; k-- {
	}

	bs[k], bs[i] = bs[i], bs[k]
	reverseBytes(bs, i+1, len(bs)-1)
	return true
}

func calcAllPermutation(s string) {
	bs := []byte(s)
	log.Println(s)
	for nextPermutation(bs) {
		log.Println(string(bs))
	}
}

func Test_CalcAllPermutation(t *testing.T) {
	calcAllPermutation("12345")
	calcAllPermutation("01")
}
