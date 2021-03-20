package microsoft

func climbStairs(n int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}
	// 考虑到每次状态转移，只会用到i, i-1, i-2三个变量，只需要两个变量记录之前值
	// 与fib序列写法一样
	prev, curr := 1, 2
	for i := 2; i < n; i++ {
		tmp := prev + curr
		prev, curr = curr, tmp
	}
	return curr
}
