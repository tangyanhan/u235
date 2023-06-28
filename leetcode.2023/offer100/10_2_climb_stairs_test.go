package offer100_test

import (
	"strconv"
	"testing"
)

var dp [101]uint64

func numWays(n int) int {
	dp[0] = 1
	dp[1] = 1
	dp[2] = 2
	for i := 3; i <= n; i++ {
		if dp[i] != 0 {
			continue
		}
		dp[i] = (dp[i-1] + dp[i-2]) % (1e9 + 7)
	}
	return int(dp[n])
}

func Test_numWays(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{2, 2},
		{7, 21},
		{0, 1},
		{92, 720754435},
		{95, 93363621},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.n), func(t *testing.T) {
			if got := numWays(tt.n); got != tt.want {
				t.Errorf("numWays() = %v, want %v", got, tt.want)
			}
		})
	}
}
