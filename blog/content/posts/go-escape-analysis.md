---
title: "Go源码分析: 逃逸分析"
date: 2020-04-05T23:25:07+08:00
tags: "Go"
---

# 什么是逃逸分析

逃逸分析(Escape Analysis)是Go在编译程序时执行的过程, 由编译器通过分析, 决定变量应当分配在栈上还是堆上.


# 在编译中进行逃逸分析

目前有代码如下:

```go
package main

import (
	"fmt"
)

type User struct {
	name string
}

func GetUsername(u *User) string {
	return u.name
}

func escapeSimple() int {
	i := 1
	j := i + 1
	return j
}

func main() {
	fmt.Println(escapeSimple())
}
```

通过在编译时增加gcflags参数, 使用类似如下命令编译:
```
go build -gcflags '-m -N -l' ./advanced/cmd/escape-analysis
```

然后获得输出如下:
```
# github.com/tangyanhan/u235/advanced/cmd/escape-analysis
advanced/cmd/escape-analysis/main.go:11:18: leaking param: u to result ~r1 level=1
advanced/cmd/escape-analysis/main.go:22:13: main ... argument does not escape
advanced/cmd/escape-analysis/main.go:22:26: escapeSimple() escapes to heap
```

这些信息,表明 GetUsername 将参数"泄露"到了返回值中, 而 escapeSimple 则逃逸到了堆中.

## 编译时参数是怎么来的?

在 go build 命令执行时, 其实包含了编译(compile), 连接(link)等多个步骤, 这里 ```-gcflags``` 实际上是传递给 ```go tool compile```的参数, 相关列表可以通过以下命令获得:

```
go tool compile --help
```

类似的, 在连接时, 通过```-ldflags``` 传递给```go tool link```, 对应参数列表, 可以通过以下命令获得:

```
go tool link --help
```

# 源码中的逃逸分析

在Go源码中, 通过注释解释了逃逸分析的运行机制. 1.14源码中, 这段注释出现在 ```src/cmd/compile/internal/gc/escape.go``` 开头.

第一段如下:

```
这里我们通过分析函数来决定Go变量应当分配到栈上, 包括那些明确调用了 new 和 make 的语句. 我们必须要保证的两点不变条件是:
1. 指向栈对象的指针不能被存在堆里
2. 指向栈对象的指针,生命周期不能超出栈对象本身(因为声明栈对象的函数在返回时已经摧毁栈帧,或者它的空间被复用于循环中的局部变量)
```

这里揭示了几点:

* 即使是明确调用 new/make 创建出来的变量, 也可能被分配到栈上
* 当函数中的变量被返回时, 它将不可能被分配到栈上. 循环中的局部变量不会被分配到堆上(一般情况)

## new/make 不一定逃逸

对于第一点, 以下面的代码为例, 就会发现 GetUsername 中通过 new 创建出的 p, 实际生命周期并没有超出函数范围. 而 return u.name, 导致参数 u 被抛出了范围.

```go
func GetUsername(u *User) string {
	p := new(User)
	p.name = "John"
	fmt.Println(p.name)
	return u.name
}
```

而分析结果也如我们所料:
```
advanced/cmd/escape-analysis/main.go:11:18: leaking param: u to result ~r1 level=1
advanced/cmd/escape-analysis/main.go:12:10: GetUsername new(User) does not escape
advanced/cmd/escape-analysis/main.go:14:13: GetUsername ... argument does not escape
```

## 循环逃逸

```go
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
```
分析结果如下:

```
advanced/cmd/escape-analysis/main.go:27:24: loop map[string]string literal does not escape
advanced/cmd/escape-analysis/main.go:33:14: loop ... argument does not escape
advanced/cmd/escape-analysis/main.go:33:14: k escapes to heap
advanced/cmd/escape-analysis/main.go:33:14: v escapes to heap
```

首先, m 虽然是个map, 但它很小, 而且 loop 自产自销, 在栈空间足够的情况下, 是可以使用的.
其次, k, v 在循环中被复用, 因此也被分配到了堆上.

### 循环逃逸带来的一个小问题

假设现有 Obj Slice, 其内部的Val如下: 1, 2, 3, 4, 5. 通过下面代码, 调用 print 后打印结果是什么?

```go
type Obj struct {
    Val int
}

func print(objs []*Obj) {
    for _, v := range objs {
        defer fmt.Println(v.Val)
        defer func() {
            fmt.Println(v.Val)
        }()
    }
}
```

## 闭包引用变量逃逸

对于生命周期, 主要就是围绕着返回值, 那么如果是闭包呢?
```go
func closure() {
	x := "hello"
	fn := func() {
		fmt.Println(x)
	}
	fn()
}
```
分析发现, x 被分配到了堆上, 闭包中引用的变量会被分配到堆上.
```
advanced/cmd/escape-analysis/main.go:20:8: closure func literal does not escape
advanced/cmd/escape-analysis/main.go:21:14: closure.func1 ... argument does not escape
advanced/cmd/escape-analysis/main.go:21:14: x escapes to heap
```

这里又引出了另一个问题, 关于闭包的实现问题... 以后再说.

## new/make 等被判定分配到栈上的阈值是多少?

我们知道, 栈的大小是有限的, 如果系统限制栈长度为8mb, 那么我们就不可能分配一个10mb的slice到栈上. 之前我们提到过有些语句, 即使我们明确使用了new/make, 创建出的对象还是可能被分配到栈上. 

那么问题来了, Go依据什么决定new/make分配到栈上呢?

1.14 ```src/cmd/compile/internal/gc/esc.go:mustHeapALloc``` 描述了这个逻辑:

```go
func mustHeapAlloc(n *Node) bool {
	// TODO(mdempsky): Cleanup this mess.
	return n.Type != nil &&
		(n.Type.Width > maxStackVarSize ||
			(n.Op == ONEW || n.Op == OPTRLIT) && n.Type.Elem().Width >= maxImplicitStackVarSize ||
			n.Op == OMAKESLICE && !isSmallMakeSlice(n))
}

// ...
var (
	// maximum size variable which we will allocate on the stack.
	// This limit is for explicit variable declarations like "var x T" or "x := ...".
	// Note: the flag smallframes can update this value.
	maxStackVarSize = int64(10 * 1024 * 1024)

	// maximum size of implicit variables that we will allocate on the stack.
	//   p := new(T)          allocating T on the stack
	//   p := &T{}            allocating T on the stack
	//   s := make([]T, n)    allocating [n]T on the stack
	//   s := []byte("...")   allocating [n]byte on the stack
	// Note: the flag smallframes can update this value.
	maxImplicitStackVarSize = int64(64 * 1024)
)
```

1. Object 自身长度不能超过栈长度
2. Object不超过最大栈变量长度(目前64位linux上是 64k)

事实上, slice/map 只是一个普通的struct, 往往实际分配都在堆上. map 在逃逸分析时只是被作为一个普通的struct, 因为其内元素大小/增长, 并不会影响其struct本身.

slice略微特殊, slice在分配和增长中有一套自己的逻辑, 如果对很小的slice也统统分配到堆上, 可能会造成大量的内存碎片. slice目前的分配阈值是64k(linux 64, Go1.14, 且未通过smallframes变更 maxImplicitStackVarSize的值). 即不超过64k, 且经过逃逸分析未逃逸的slice, 会被分配到栈上, 而不是堆上.

# 相关代码

相关代码放在我的github仓库中: [advanced/cmd/escape-analysis/main.go](https://github.com/tangyanhan/u235/blob/master/go/advanced/cmd/escape-analysis/main.go)