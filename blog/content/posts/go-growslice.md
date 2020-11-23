---
title: "Go源码分析: Slice内存增长方式"
date: 2020-04-04T22:41:18+08:00
tags: "Go"
---

源文件: runtime/slice.go

slice 是一个连续的内存结构, 由3部分组成: 内存区域, 长度 和 可用长度(cap):

```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

# make

创建一个slice， 实际上有3个参数： ```make([]type, len, cap)```

当省略参数cap时，默认cap=len

cap决定了实际分配内存区域的大小， 在进行 append/copy等操作时， 如果长度没有超过cap， 则不需要重新分配内存， 也不需要在内存中移动数据

make slice原型： func makeslice

```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return mallocgc(mem, et, true)
}
```

# append

当调用 append 向 slice 添加元素时, 会发生什么?

在所需新的cap低于1024时, 每次slice达到极限, 需要扩充时, 都会乘以原来长度的两倍.

当超过1024时, 每次都会增长原长度的 1/4, 直至新的长度 > 预期长度 为止.

```go
	newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if old.len < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
```

验证扩充长度变化:

```go
func TestAppend(t *testing.T) {
	var s []int

	t.Log(cap(s))
	for i := 0; i < 1025; i++ {
		s = append(s, i)
		t.Log(cap(s))
	}
}
```

此外, 每次扩张都会需要内存迁移 memmove 进行迁移.

因此, 当slice发生扩容时, 实际包含了:

1. 计算新的长度
2. 分配内存
3. memmove

在进行 memmove 时, 不会对原内存区域进行清零. 因此, 假如原 slice 包含敏感信息, 此时对内存进行dump, 是可以找到残留信息的 (密码学实现需要关注这一点).

另一个问题, 分配内存时是否一定会通过系统调用分配内存? 需要在单独研究中分析.

# 启示

从以上分析中, 我们知道 growslice 是一个代价较高的过程. 如何避免或者减少这种影响呢?

答案就是在分配时尽可能预设合适的cap. 如无必要, 避免在不预分配cap的情况下不断append.

## append slice

以下代码, 源slice一经确定, 但 s 最终 cap为16, 浪费了空间, 且发生了多次调整: 0-> 1 -> 2 -> 4 -> 8 -> 16

```go
func TestAppend(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var s []int

	t.Log(cap(s))
	for _, v := range src {
		s = append(s, v)
		t.Log(cap(s))
	}
}
```

## Buffer 缓存

bytes.Buffer 一般认为比直接将字符串相加效率更高, 但其存储使用的 byte slice, 也需要同样的扩容过程. 因此较好的策略是, 初始化时给它一个适度大小的slice作为起始内存区域.

```go
    // 在某段函数中, 尝试计算和预测所需buf长度
	totalLen += len(l)*4 + 1
	buf := make([]byte, totalLen)
	b := bytes.NewBuffer(buf)
```
