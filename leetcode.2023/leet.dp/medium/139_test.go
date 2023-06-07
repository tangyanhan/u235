package medium

import "testing"

/*
本质上是算法导论cutRod的简化版，s就是那根棍子，选择不同的单词，如果基于该单词长度的下一个位置已经是匹配的，那么就可以进行接下来的匹配
*/
func wordBreak(s string, wordDict []string) bool {
	dp := make([]bool, len(s)+1)
	dp[len(s)] = true

	for i := len(s) - 1; i >= 0; i-- {
		for _, word := range wordDict {
			nextIdx := i + len(word)
			if nextIdx <= len(s) && dp[nextIdx] {
				// Compare if word match
				match := true
				for j := 0; j < len(word); j++ {
					if word[j] != s[i+j] {
						match = false
						break
					}
				}
				if match {
					dp[i] = true
					break
				}
			}
		}
	}

	return dp[0]
}

func Test_wordBreak(t *testing.T) {
	tests := []struct {
		name     string
		wordDict []string
		want     bool
	}{
		{"leetcode", []string{"leet", "code"}, true},
		{"applepenapple", []string{"apple", "pen"}, true},
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wordBreak(tt.name, tt.wordDict); got != tt.want {
				t.Errorf("wordBreak() = %v, want %v", got, tt.want)
			}
		})
	}
}
