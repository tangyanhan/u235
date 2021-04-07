package casual

import (
	"fmt"
)

// Prime is a randomly chosen number
const Prime = uint64(37)

// HashStr 将字符串映射到一个整型上
func HashStr(s string) (hash uint64, pow uint64) {
	hash = uint64(0)
	for i := 0; i < len(s); i++ {
		hash = hash*Prime + uint64(s[i])
		fmt.Println("Hash=", hash)
	}
	pow = 1
	sq := Prime
	for i := len(s); i > 0; i >>= 1 {
		if i&1 != 0 {
			pow *= sq
			fmt.Println("pow=", pow, "i=", i)
		}
		sq *= sq
		fmt.Println("sq=", sq, "i=", i)
	}
	return
}

// IndexString first occurrence of substr
func IndexString(s, substr string) int {
	if len(s) <= len(substr) {
		if s == substr {
			return 0
		}
		return -1
	}
	hashss, pow := HashStr(substr)
	fmt.Println("substr hash=", hashss)
	n := len(substr)
	var h uint64 = 0
	for i := 0; i < n; i++ {
		h = h*Prime + uint64(s[i])
	}
	fmt.Println("Source hash=", h)
	if h == hashss && s[0:n] == substr {
		return 0
	}
	for i := n; i < len(s); {
		h *= Prime
		h += uint64(s[i])
		h -= pow * uint64(s[i-n])
		i++
		fmt.Println("str=", s[i-n:i], "hash=", h)
		if h == hashss && s[i-n:i] == substr {
			return i - n
		}
	}
	return -1
}
