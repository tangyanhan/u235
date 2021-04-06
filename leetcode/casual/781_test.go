package casual

import (
	"testing"
)

func numRabbits(answers []int) int {
	numMap := make(map[int]int)
	for _, v := range answers {
		numMap[v+1]++
	}
	var ans int
	for num, cnt := range numMap {
		if cnt <= num {
			ans += num
		} else {
			for cnt > 0 {
				ans += num
				cnt -= num
			}
		}
	}
	return ans
}

func Test_numRabbits(t *testing.T) {
	tests := []struct {
		name    string
		answers []int
		want    int
	}{
		{
			answers: []int{},
			want:    0,
		},
		{
			answers: []int{1, 1, 2},
			want:    5,
		},
		{
			answers: []int{10, 10, 10},
			want:    11,
		},
		{
			answers: []int{1, 0, 1, 0, 0},
			want:    5,
		},
		{
			answers: []int{0, 0, 1, 1, 1},
			want:    6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numRabbits(tt.answers); got != tt.want {
				t.Errorf("numRabbits() = %v, want %v", got, tt.want)
			}
		})
	}
}
