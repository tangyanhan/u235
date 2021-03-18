package microsoft

import (
	"math"
	"sort"
	"testing"
)

func coinChange(coins []int, amount int) int {
	// TODO: fix this
	sort.Ints(coins)
	var changeRemainAmount func(int, int, uint32)
	ans := uint32(math.MaxUint32)
	changeRemainAmount = func(curCoin int, amount int, curCnt uint32) {
		if amount == 0 {
			if curCnt < ans {
				ans = curCnt
			}
			return
		}
		if curCoin < 0 {
			return
		}
		coinValue := coins[curCoin]
		maxNum := amount / coinValue
		totalValue := maxNum * coinValue
		for i := maxNum; maxNum >= 0 && curCnt+uint32(maxNum) < ans; maxNum-- {
			changeRemainAmount(curCoin-1, amount-totalValue, curCnt+uint32(i))
			totalValue -= coinValue
		}
	}
	changeRemainAmount(len(coins)-1, amount, 0)
	if ans == math.MaxUint32 {
		return -1
	}
	return int(ans)
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
		{
			args: args{
				coins:  []int{1, 2, 5},
				amount: 11,
			},
			want: 3,
		},
		{
			args: args{
				coins:  []int{1},
				amount: math.MaxInt32,
			},
			want: math.MaxInt32,
		},
		{
			args: args{
				coins:  []int{186, 419, 83, 408},
				amount: 6249,
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coinChange(tt.args.coins, tt.args.amount); got != tt.want {
				t.Errorf("coinChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
