package microsoft

import (
	"reflect"
	"testing"
)

var dialMap = []string{
	"",
	"",
	"abc",
	"def",
	"ghi",
	"jkl",
	"mno",
	"pqrs",
	"tuv",
	"wxyz",
}

func letterCombinations(digits string) []string {
	if len(digits) == 0 {
		return nil
	}
	alts := make([]string, len(digits))
	var totalNum int
	for i, c := range digits {
		alts[i] = dialMap[int(c)-'0']
		totalNum *= len(alts[i])
	}
	results := make([]string, 0, totalNum)
	value := make([]int, len(alts))
	strBuf := make([]byte, len(digits))
	for {
		for i := 0; i < len(value); i++ {
			strBuf[i] = alts[i][value[i]]
		}
		results = append(results, string(strBuf))

		var carry int
		value[len(value)-1]++
		for i := len(alts) - 1; i >= 0; i-- {
			value[i] += carry
			carry = 0
			if value[i] >= len(alts[i]) {
				carry = 1
				value[i] = 0
			}
		}
		if carry != 0 {
			break
		}
	}

	return results
}

func Test_letterCombinations(t *testing.T) {
	type args struct {
		digits string
	}
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			in:   "2",
			want: []string{"a", "b", "c"},
		},
		{
			in:   "23",
			want: []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := letterCombinations(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("letterCombinations() = %v, want %v", got, tt.want)
			}
		})
	}
}
