package microsoft

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func multiply(num1 string, num2 string) string {
	if num1 == "0" || num2 == "0" {
		return "0"
	}
	var idx int
	buf := make([]byte, len(num1)+len(num2))
	for ; idx < len(num2); idx++ {
		var carry byte
		mul := num2[len(num2)-1-idx] - '0'
		for i := len(num1) - 1; i >= 0; i-- {
			v := (num1[i]-'0')*mul + carry
			writeAt := len(buf) - 1 - idx - (len(num1) - 1 - i)
			sum := buf[writeAt]
			sum += v
			carry = sum / 10
			sum %= 10
			buf[writeAt] = sum
		}
		if carry != 0 {
			buf[len(buf)-1-idx-len(num1)] += carry
		}
	}
	begin := -1
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 && begin == -1 {
			begin = i
		}
		if begin != -1 {
			buf[i] += '0'
		}
	}

	return string(buf[begin:])
}

func Test_multiply(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		a := rand.Int() % 100000
		b := rand.Int() % 100000
		got := multiply(strconv.Itoa(a), strconv.Itoa(b))
		c, err := strconv.Atoi(got)
		if err != nil {
			t.Fatal(err)
		}
		if c != a*b {
			t.Fatalf("#%d failed: %d * %d = %d, got=%s", i, a, b, a*b, got)
		}
	}
}
