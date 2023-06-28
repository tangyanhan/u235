package offer100_test

import (
	"math"
	"testing"
)

func maxProfit(prices []int) int {
	minValue := math.MaxInt
	maxProfit := 0
	for _, p := range prices {
		if p < minValue {
			minValue = p
			continue
		}
		profit := p - minValue
		if profit > maxProfit {
			maxProfit = profit
		}
	}
	return maxProfit
}

func Test_maxProfit(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"general profit", []int{7, 1, 5, 3, 6, 4}, 5},
		{"no profit", []int{7, 6, 4, 3, 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxProfit(tt.prices); got != tt.want {
				t.Errorf("maxProfit() = %v, want %v", got, tt.want)
			}
		})
	}
}
