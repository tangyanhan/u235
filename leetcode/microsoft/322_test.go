package microsoft

import (
	"testing"
)

func coinChange(coins []int, amount int) int {
	// 状态： 金额i，使用coins拼组的最小数量dp[i]
	// 转移方程： dp[i] = min( dp[i], min(dp[i-c0], ..., dp[i-cj])+1)
	// dp[0] = 0
	// dp[i-cj]+1，表示当金额达到i时，从i中扣掉金币cj，正好得到已经计算过的金额dp[i-cj]，这表示只要我们从该值加上一枚硬币cj，即可到达值i
	// 只要我们找到所有“扣掉硬币j”+1的最小值，就得到了dp[i]
	// 预先填充dp[0]以外的所有元素为10001(amount+1)，这样如果跳到某个不可能到达的值，就会计算出一个比amount更高的数值
	dp := make([]int, amount+1)
	for i := 1; i < len(dp); i++ {
		dp[i] = 10001
	}
	for i := 1; i <= amount; i++ {
		for j := 0; j < len(coins); j++ {
			if coins[j] <= i {
				v := dp[i-coins[j]] + 1
				if v < dp[i] {
					dp[i] = v
				}
			}
		}
	}

	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
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
				amount: 10000,
			},
			want: 10000,
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
