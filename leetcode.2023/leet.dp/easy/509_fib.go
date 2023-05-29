package easy

var fibMap = map[int]int{
	0: 0,
	1: 1,
}

func fibMemo(n int) int {
	if v, ok := fibMap[n]; ok {
		return v
	}

	ret := fibMemo(n-1) + fibMemo(n-2)
	fibMap[n] = ret
	return ret
}

var fibArr = [31]int{0, 1}

func fibTable(n int) int {
	ret := fibArr[n]
	if ret != 0 || n == 0 {
		return ret
	}

	ret = fibTable(n-1) + fibTable(n-2)
	fibArr[n] = ret
	return ret
}

func fib(n int) int {
	if n < 2 {
		return n
	}

	var ret int
	n1, n2 := 0, 1

	for i := 2; i <= n; i++ {
		ret = n1 + n2
		n1, n2 = n2, ret
	}
	return ret
}
