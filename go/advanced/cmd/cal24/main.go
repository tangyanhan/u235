package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const N = 4

type Num struct {
	Value float64
	Expr  string
}

func (n Num) String() string {
	if n.Expr[0] == '(' {
		return string(n.Expr[1 : len(n.Expr)-1])
	}
	return n.Expr
}

type NumSet map[Num]bool

func (f *NumSet) Union(b NumSet) {
	for k, v := range b {
		(*f)[k] = v
	}
}

func Equal(a, b float64) bool {
	v := a - b
	if v < 0.0 {
		v = -v
	}
	return v < 0.00000001
}

func Fork(nums []int, sets []NumSet, m int) NumSet {
	if len(sets[m]) != 0 {
		return sets[m]
	}
	for i := 1; i < m/2; i++ {
		if (i & m) == i {
			s1 := Fork(nums, sets, i)
			s2 := Fork(nums, sets, m-i)
			for a, _ := range s1 {
				for b, _ := range s2 {
					sets[m][Num{
						Value: a.Value + b.Value,
						Expr:  "(" + a.Expr + "+" + b.Expr + ")",
					}] = true
					sets[m][Num{
						Value: a.Value - b.Value,
						Expr:  "(" + a.Expr + "-" + b.Expr + ")",
					}] = true
					sets[m][Num{
						Value: b.Value - a.Value,
						Expr:  "(" + b.Expr + "-" + a.Expr + ")",
					}] = true
					sets[m][Num{
						Value: a.Value * b.Value,
						Expr:  "(" + a.Expr + "*" + b.Expr + ")",
					}] = true
					if !Equal(b.Value, 0.0) {
						sets[m][Num{
							Value: a.Value / b.Value,
							Expr:  "(" + a.Expr + "/" + b.Expr + ")",
						}] = true
					}
					if !Equal(a.Value, 0.0) {
						sets[m][Num{
							Value: b.Value / a.Value,
							Expr:  "(" + b.Expr + "/" + a.Expr + ")",
						}] = true
					}
				}
			}
		}
	}
	return sets[m]
}

func main() {
	area := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	testN := 100
	var nums [N]int
	var sets [1 << N]NumSet
	for t := 0; t < testN; t++ {
		// shuffle
		for i := 0; i < len(area); i++ {
			x := rand.Intn(len(area))
			area[i], area[x] = area[x], area[i]
		}
		nums[0] = area[0]
		nums[1] = area[1]
		nums[2] = area[2]
		nums[3] = area[3]

		for i := 0; i < (1 << N); i++ {
			sets[i] = make(NumSet)
		}
		for i := 0; i < N; i++ {
			sets[1<<i][Num{
				Value: float64(nums[i]),
				Expr:  strconv.Itoa(nums[i]),
			}] = true
		}
		// Start from set of 4 elements
		Fork(nums[:], sets[:], (1<<N)-1)
		fmt.Println("Numbers:", nums)
		for num, _ := range sets[(1<<N)-1] {
			if Equal(num.Value, 24.0) {
				fmt.Println(num)
			}
		}
	}
}
