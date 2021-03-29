package pingcap

import (
	"math"
	"math/rand"
	"testing"
)

func rand7() int {
	return rand.Intn(7) + 1
}

// 文章基于这样一个事实 (randX() - 1)*Y + randY() 可以等概率的生成[1, X * Y]范围的随机数

func rand10() int {
	var sum int
	for i := 0; i < 10; i++ {
		sum += rand7()
	}
	return sum%10 + 1
}

func TestRandDistribution(t *testing.T) {
	m := make(map[int]int)
	const N = 1000
	const Max = 10
	const delta = N / 20
	for i := 0; i < N; i++ {
		m[rand10()]++
	}
	ideal := N / Max
	for k, cnt := range m {
		diff := int(math.Abs(float64(cnt) - float64(ideal)))
		if diff > delta {
			t.Fatalf("Bad distribution for value %d: expect %d, got %d", k, ideal, cnt)
		}
	}
}
