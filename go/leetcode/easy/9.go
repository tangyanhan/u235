package easy

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	p := make([]int, 0)
	for x > 0 {
		d := x % 10
		x /= 10
		p = append(p, d)
	}

	i := 0
	j := len(p) - 1
	for i < j {
		if p[i] != p[j] {
			return false
		}
		i++
		j--
	}
	return true
}
