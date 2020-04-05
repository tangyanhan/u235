package main

import (
	"fmt"
)

type User struct {
	name string
}

type Obj struct {
	Val int
}

func GetUsername(u *User) string {
	p := new(User)
	p.name = "John"
	fmt.Println(p.name)
	return u.name
}

func closure() {
	x := "hello"
	fn := func() {
		fmt.Println(x)
	}
	fn()
}

func loop() {
	m := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func escapeSimple() int {
	i := 1
	j := i + 1
	return j
}

func printObj(objs []*Obj) {
	for _, v := range objs {
		defer fmt.Println("a=", v.Val)
		defer func() {
			fmt.Println("b=", v.Val)
		}()
	}
}

// const Ki = 1024
// const Mi = 1024 * Ki

// func largeObj() {
// 	large := make([]byte, 64*Ki-1)
// 	for i := 0; i < cap(large); i++ {
// 		large[i] = byte(1)
// 	}
// }

type Array [1]int64

func largeArray() {
	var arr Array
	fmt.Println(arr)
}

func main() {
	fmt.Println(escapeSimple())
	objs := []*Obj{
		&Obj{1}, &Obj{2}, &Obj{3}, &Obj{4}, &Obj{5},
	}
	printObj(objs)
	largeArray()
}
