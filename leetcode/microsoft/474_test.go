package microsoft

import (
	"fmt"
	"testing"
)

func printDP(dp [][]int, prefix string) {
	fmt.Println("===Current DP:", prefix)
	for _, row := range dp {
		fmt.Println(row)
	}
}

func findMaxForm(strs []string, m int, n int) int {
	// TODO: the original solution is not correct
	return 0
}

// 表面上看这道题能在LeetCode过，但无法通过下面的测试用例
func findMaxFormBad(strs []string, m int, n int) int {
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for _, s := range strs {
		var zeroNum, oneNum int
		for _, c := range s {
			if c == '0' {
				zeroNum++
			} else {
				oneNum++
			}
		}
		for i := m; i >= zeroNum; i-- {
			for j := n; j >= oneNum; j-- {
				useThisStr := 1 + dp[i-zeroNum][j-oneNum]
				if useThisStr > dp[i][j] {
					dp[i][j] = useThisStr
				}
				printDP(dp, fmt.Sprint("i=", i, "j=", j, "zeroNum=", zeroNum, "oneNum=", oneNum))
			}
		}
	}
	return dp[m][n]
}

func Test_findMaxForm(t *testing.T) {
	type args struct {
		strs []string
		m    int
		n    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				strs: []string{"10, 0001", "111001", "1", "0"},
				m:    5,
				n:    3,
			},
			want: 4,
		},
		{
			args: args{
				strs: []string{"10", "0", "1"},
				m:    1,
				n:    1,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMaxForm(tt.args.strs, tt.args.m, tt.args.n); got != tt.want {
				t.Errorf("findMaxForm() = %v, want %v", got, tt.want)
			}
		})
	}
}
