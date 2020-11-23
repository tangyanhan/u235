---
title: "Cas与Mutex"
date: 2020-04-26T14:56:49+08:00
tag: "Go"
draft: true
---

# CAS

CAS (Compare And Swap) 是一个原子性操作, 通过CPU LOCK指令前缀, 避免程序运行于多核处理器时, 每个核心各自维护的缓存并对其读写导致数据不一致的问题.

Go提供了 ```atomic``` 包封装了一系列操作, 因为处理器架构不同, 指令也有所不同. 对于AMD64架构, 其封装的Cas操作封装于 ```runtime/internal/atomic/asm_amd64.s``` 中:

```asm
// bool Cas(int32 *val, int32 old, int32 new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	} else
//		return 0;
TEXT runtime∕internal∕atomic·Cas(SB),NOSPLIT,$0-17
    // 将 val 放入 BX
	MOVQ	ptr+0(FP), BX
    // 将 old 放入 AX
	MOVL	old+8(FP), AX
    // 将 new 放入 CX
	MOVL	new+12(FP), CX
    // LOCK 前缀用于处理器访问缓存或内存的锁定, CMPXCHGL 指令默认不加也是可以的
	LOCK 
    // CMPXCHGL  src, dst
    // 比较 AX 与 dst, 如果 dst == AX, 则将src赋值给 dst
    // 如果成功, 则ZF=1, 否则ZF=0
	CMPXCHGL	CX, 0(BX)
	SETEQ	ret+16(FP)
	RET
```

# sync.Mutex

## 1. Mutex不可复制

首先看 sync/mutex.go 中的实现, 可以看到其基本结构就是两个整数. 因此, 当把一个mutex通过函数参数等形式复制时, 新的mutex继承了原mutex的状态, 但是两者已经不是同一个mutex了, 对其中一个Unlock不会影响到另一个. 因此mutex不可复制.

```
// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
type Mutex struct {
	state int32
	sema  uint32
}
```

## 2. Lock/Unlock

在简单情况下, 只需要将 state 变量从 0 变为1即可完成Lock, 从1变回0即可完成Unlock.

当有多个goroutine尝试加锁/解锁时, 如果仍旧沿用这个简单机制, 那么 goroutine 获得锁的顺序是不确定的. 我们能否让尝试加锁的goroutine分出一个先后顺序, 防止后来的加锁请求反而先于先来的加锁成功呢?

