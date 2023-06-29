package offer100_test

import (
	"fmt"
	"strconv"
	"testing"
)

func pickMin(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= a && b <= c {
		return b
	}

	return c
}

func nthUglyNumber(n int) int {
	var nums [1691]int
	m2, m3, m5 := 0, 0, 0
	x2, x3, x5 := 2, 3, 5

	nums[0] = 1
	for i := 1; i < n; i++ {
		minValue := pickMin(x2, x3, x5)
		nums[i] = minValue
		fmt.Println("Alternate:", x2, x3, x5, "Index=", m2, m3, m5, "min=", minValue)

		if minValue == x2 {
			m2++
			x2 = nums[m2] * 2
		}
		if minValue == x3 {
			m3++
			x3 = nums[m3] * 3
		}
		if minValue == x5 {
			m5++
			x5 = nums[m5] * 5
		}
	}

	return nums[n-1]
}

func Test_NthUglyNumber(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12}

	for i, expect := range data {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			got := nthUglyNumber(i + 1)
			if got != expect {
				t.Fatalf("%d th Expect:%d Got:%d", i+1, expect, got)
			}
		})
	}
}
