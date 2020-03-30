package easy

func isValid(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	stack := make([]rune, len(s))
	j := -1
	pop := func(expect rune) bool {
		if j < 0 {
			return false
		}
		if stack[j] != expect {
			return false
		}
		j--
		return true
	}
	for _, r := range s {
		switch r {
		case '{', '[', '(':
			j++
			stack[j] = r
		case '}':
			if !pop('{') {
				return false
			}
		case ']':
			if !pop('[') {
				return false
			}
		case ')':
			if !pop('(') {
				return false
			}
		}
	}

	return j == -1
}
