package microsoft

import (
	"fmt"
	"testing"
)

func TestWalkArray(t *testing.T) {
	var m [4][4]int
	num := 1
	for delta := 0; delta < len(m); delta++ {
		i, j := 0, delta
		for i < len(m) && j < len(m) {
			m[i][j] = num
			i++
			j++
		}
		num++
	}
	for _, row := range m {
		for _, v := range row {
			fmt.Printf("%x ", v)
		}
		fmt.Println("")
	}
}

func maxLength(arr []string) int {
	dp := make([][]int32, len(arr))
	for i := range dp {
		dp[i] = make([]int32, len(arr))
	}
	var maxLen int
	for i, s := range arr {
		if len(s) == 26 {
			return 26
		}
		if len(s) > maxLen {
			maxLen = len(s)
		}
		for j := range s {
			dp[i][i] |= (1 << int32(s[j]-'a'))
		}
	}
	for _, row := range dp {
		for _, v := range row {
			fmt.Printf("%x ", v)
		}
		fmt.Println("")
	}

	getLength := func(v int32) int {
		var length int
		for i := 0; i < 26; i++ {
			if v&(1<<i) != 0 {
				length++
			}
		}
		return length
	}

	for i := 0; i < len(dp); i++ {
		for j := 1; j < len(dp[i]); j++ {

		}
	}

	for i := 0; i < len(dp); i++ {
		for j := i + 1; j < len(dp[i]); j++ {
			// test if has joins
			if dp[j][j]&dp[i][i] == 0 {
				dp[i][j] = dp[j][j] | dp[i][i]
				length := getLength(dp[i][j])
				if length > maxLen {
					maxLen = length
				}
				if length == 26 {
					return 26
				}
			} else {
				dp[i][j] = 0
			}
		}
	}
	return maxLen
}

func Test_maxLength(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want int
	}{
		{
			in:   []string{"cha", "r", "act", "ers"},
			want: 6,
		},
		{
			in:   []string{"a", "b", "bc", "bd", "ef", "ac", "xyz"},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxLength(tt.in); got != tt.want {
				t.Errorf("maxLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
