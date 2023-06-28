package offer100_test

import (
	"log"
	"strconv"
	"testing"
)

func translateNum(num int) int {
	if num == 0 {
		return 1
	}

	parts := make([]int, 0, 65)

	for num > 0 {
		n := num % 10
		num /= 10
		parts = append(parts, n)
	}

	log.Println(parts)

	dp := make([]int, len(parts))
	dp[0] = 1
	for i := 1; i < len(dp); i++ {
		val := parts[i]*10 + parts[i-1]
		if val < 26 && val >= 10 {
			if i-2 >= 0 {
				dp[i] = dp[i-1] + dp[i-2]
			} else {
				dp[i] = dp[i-1] + 1
			}

		} else {
			dp[i] = dp[i-1]
		}
	}

	log.Println(dp)

	return dp[len(dp)-1]
}

func Test_translateNum(t *testing.T) {
	tests := []struct {
		num  int
		want int
	}{
		{0, 1},
		{12258, 5},
		{25, 2},
		{1068385902, 2},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.num), func(t *testing.T) {
			if got := translateNum(tt.num); got != tt.want {
				t.Errorf("translateNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
