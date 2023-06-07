package hard

import (
	"log"
	"testing"
)

/*
最小问题：当t只剩1个字母时，只需要统计该字母在s中的出现频率即可
扩增问题：当t增加到两个字母时，最右侧字母组合是重复问题
状态转移：取s位置与t位置两个因素为变量，如果当前s位置与t位置字母不同，则不可能，取更短状态下的值。如果相同，取之前计算的s+1/t+1
状态终止：t已经移动到位置0
*/
func numDistinct(s string, t string) int {
	lenS, lenT := len(s), len(t)
	if lenT > lenS {
		return 0
	}

	dp := make([][]int, lenS+1)
	dp[lenS] = make([]int, lenT+1)
	dp[lenS][lenT] = 1
	for is := lenS - 1; is >= 0; is-- {
		dp[is] = make([]int, lenT+1)
		dp[is][lenT] = 1
		for it := lenT - 1; it >= 0; it-- {
			if t[it] == s[is] {
				dp[is][it] = dp[is+1][it+1] + dp[is+1][it]
			} else {
				dp[is][it] = dp[is+1][it]
			}
		}
		log.Println("is=", is, string(s[is:]))
		log.Println(dp[is])
	}

	return dp[0][0]
}

func numDistinctMem(s string, t string) int {
	lenS, lenT := len(s), len(t)
	if lenT > lenS {
		return 0
	}

	prevRow := make([]int, lenT+1)
	prevRow[lenT] = 1
	curRow := make([]int, lenT+1)
	curRow[lenT] = 1
	for is := lenS - 1; is >= 0; is-- {
		for it := lenT - 1; it >= 0; it-- {
			if t[it] == s[is] {
				curRow[it] = prevRow[it+1] + prevRow[it]
			} else {
				curRow[it] = prevRow[it]
			}
		}
		prevRow, curRow = curRow, prevRow
	}

	return prevRow[0]
}

func Test_numDistinct(t *testing.T) {
	type args struct {
		s string
		t string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"ttt", args{"ttat", "t"}, 3},
		{"aaattt", args{"aaattt", "at"}, 9},
		{"rabbit", args{"rabbbit", "rabbit"}, 3},
		{"bag", args{"babgbag", "bag"}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numDistinct(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("numDistinct() = %v, want %v", got, tt.want)
			}
		})
	}
}
