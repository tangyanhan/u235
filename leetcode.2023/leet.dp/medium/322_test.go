package medium_test

import (
	"math"
	"sort"
	"testing"
)

func coinChangeDPMap(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}

	coinMap := make(map[int]int)
	coinMap[0] = 0
	for _, coin := range coins {
		coinMap[coin] = 1
	}

	sort.Ints(coins)
	for i := 1; i <= amount; i++ {
		num := math.MaxInt
		for _, coin := range coins {
			if coin > i {
				break
			}

			remain := i - coin
			n, ok := coinMap[remain]
			if !ok {
				continue
			}

			if num > n+1 {
				num = n + 1
			}
		}

		if num != math.MaxInt {
			coinMap[i] = num
		}
	}

	if v, ok := coinMap[amount]; ok {
		return v
	}

	return -1
}

func coinChange(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}

	dp := make([]int, amount+1)

	const unDone = -2
	const invalid = -1

	for i := 0; i < len(dp); i++ {
		dp[i] = unDone
	}

	dp[0] = 0
	sort.Ints(coins)
	for _, value := range coins {
		if value > amount {
			break
		}
		dp[value] = 1
	}

	for i := 1; i <= amount; i++ {
		num := math.MaxInt
		for _, value := range coins {
			if value > i {
				break
			}

			remain := i - value
			if dp[remain] == invalid {
				continue
			}

			if dp[remain]+1 < num {
				num = dp[remain] + 1
			}
		}
		if num == math.MaxInt {
			dp[i] = invalid
		} else {
			dp[i] = num
		}
	}

	return dp[amount]
}

func coinChangeMemo(coins []int, amount int) int {
	if amount == 0 {
		return 0
	}

	coinMap := make(map[int]int)

	for _, value := range coins {
		coinMap[value] = 1
	}

	sort.Ints(coins)

	var changeFn func(int) int

	changeFn = func(amount int) int {
		if amount == 0 {
			return 0
		}

		num, ok := coinMap[amount]
		if ok {
			return num
		}

		num = math.MaxInt
		for _, coin := range coins {
			if amount < coin {
				break
			}

			n := changeFn(amount-coin) + 1
			if n == 0 { // Invalid
				continue
			}
			if num > n {
				num = n
			}
		}

		if num == math.MaxInt {
			coinMap[amount] = -1
			return -1
		}

		coinMap[amount] = num
		return num
	}

	return changeFn(amount)
}

func Test_coinChange(t *testing.T) {
	type args struct {
		coins  []int
		amount int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Doable", args{coins: []int{1, 2, 5}, amount: 11}, 3},
		{"Invalid", args{coins: []int{2}, amount: 3}, -1},
		{"Zero", args{coins: []int{1}, amount: 0}, 0},
		{"Large", args{coins: []int{186, 419, 83, 408}, amount: 6249}, 20},
		{"Big Coin", args{coins: []int{186, 419, 83, 408, 10000, 10001}, amount: 6249}, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coinChange(tt.args.coins, tt.args.amount); got != tt.want {
				t.Errorf("coinChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
