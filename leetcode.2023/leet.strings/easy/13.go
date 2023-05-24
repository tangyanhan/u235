package easy

var romanMap map[rune]int = map[rune]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

func romanToInt(s string) int {
	nums := make([]int, len(s))
	for i, c := range s {
		nums[i] = romanMap[c]
	}

	var result int
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] < nums[i+1] {
			result -= nums[i]
		} else {
			result += nums[i]
		}
	}
	return result + nums[len(nums)-1]
}

type RomanRune struct {
	c rune
	v uint16
}

var romanArray []RomanRune = []RomanRune{
	{'M', 1000},
	{'D', 500},
	{'C', 100},
	{'L', 50},
	{'X', 10},
	{'V', 5},
	{'I', 1},
}

func getRomanValue(c rune) (int, int) {
	for offset, v := range romanArray {
		if v.c == c {
			return offset, int(v.v)
		}
	}
	panic("invalid value")
}

func romanToIntMem(s string) int {
	var result int
	var lastOffset int
	var lastValue int
	for _, c := range s {
		offset, v := getRomanValue(c)
		if lastOffset > offset {
			lastValue = -lastValue
		}
		result += lastValue
		lastValue = v
		lastOffset = offset
	}
	result += lastValue

	return result
}
