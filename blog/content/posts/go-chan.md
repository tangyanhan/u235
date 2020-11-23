---
title: "Go源码分析: chan 实现"
date: 2020-05-13T11:16:24+08:00
tags: "Go"
katex: true
---

# chan组成

chan 可以简单定义为一个"带锁的队列", 当make(chan type), 不带数量时, 队列最大长度为1.

chan的实现需要达成以下目标:

1. 读阻塞
2. 写阻塞
3. 可读时, 使阻塞的goroutine恢复
4. 可写时, 使阻塞的goroutine恢复

chan的主要实现在```runtime/chan.go```中, 其主要构成如图:

![chan](/images/post/go-chan.png)

当 make(chan int, 4) 时, buf 能否分配到4个int长度的内存. 通过 sendx/recvx 标记读写的位置, 然后读/写过程分别加锁完成.

当一个 chan 可读/可写时, 此时可能有多个 goroutine 都在等待, 需要以FIFO顺序释放一个被阻塞的goroutine, 因此hchan结构中包含了sendq和recvq, 分别记录被阻塞的写goroutine 以及被阻塞的读 goroutine.


# 语法糖

与go fn()启动goroutine一样, chan的很多使用, 本质上是编译器提供的语法糖. 下面提供一些与源码的对照:

## 从chan中读取值

```go
v := <- c

// Expand by compiler:
runtime.chanrecv1(&c, &v)
```

## 读取值的同时知道chan是否关闭

```go
v, ok := <- c

// Expand by compiler:
ok := runtime.chanrecv2(&c, &v)
```

## 写入chan

```go
c <- v

// Expand:
chansend1(&c, &v)
```

## 非阻塞写入

```go
c := make(chan int)
select {
    case c <- v:
        // ... foo
    default:
        // ... bar
}

// Expand:
if selectnbsend(&c, v) {
    // ... foo
} else {
    // ... bar
}
```

## 非阻塞读取

```go
select {
    case v := <- c
    // ... foo
    default:
    // ... bar
}

// Expand:
if selectnbrecv(&c, v) {
    // ... foo
} else {
    // ... bar
}
```

# select

这个实现并不在```chan.go```中. 写下面一段程序:
```go
	a := make(chan struct{})
	c := make(chan struct{})
	select {
	case <-a:
		state += 1
	case <-c:
		state += 2
	default:
		state += 9
	}
```

通过```go tool objdump -S mychan > mychan.s``` 获得汇编与源码的对照, 发现使用的是```runtime.selectgo```.

源码在 ```runtime/select.go```中.

关于select, 我们需要解答以下问题:

## 在select中读写nil chan

当select的分支中涉及为nil chan时, select 会如何读取它们? 如果是写呢?

在 selectgo 入口处, 对nil chan进行了处理, 将其替换为一个空的scase{}:

```go
	// Replace send/receive cases involving nil channels with
	// caseNil so logic below can assume non-nil channel.
	for i := range scases {
		cas := &scases[i]
		if cas.c == nil && cas.kind != caseDefault {
			*cas = scase{}
		}
    }
```

空的```scase.kind == caseNil```, 搜索 caseNil, 发现相关代码都是这样的:

```go
// 首次读取(处理default), 忽略
loop:
// ...
switch cas.kind {
case caseNil:
    continue
// goroutine入队列, 忽略
if cas.kind == caseNil {
    continue
}
// goroutine出队列, 忽略
if k.kind == caseNil {
    continue
}

```

即select会忽略nil chan. 基于以上内容, 以下代码的输出是什么?

```go
	var invalid chan int
	valid := make(chan int)

	go func(x int) {
		time.Sleep(time.Second)
		valid <- x
	}(2333)
	var v int
	select {
	case v = <-valid:
	case v = <-invalid:
	}

    fmt.Println(v)
```

## 第一步: 尝试读写,处理default

select首先会尝试遍历一遍所有的分支, 看是否有碰巧已经可读写, 或者是直接到 default 分支返回.

遍历顺序(pollorder) 是随机的, 通过洗牌进行随机交换打乱顺序:

```go
// generate permuted order
for i := 1; i < ncases; i++ {
    j := fastrandn(uint32(i + 1))
    pollorder[i] = pollorder[j]
    pollorder[j] = uint16(i)
}
```

因此, 在多分支select带有 default 的情况下, 必定会在第一遍返回:

```go
    case caseDefault:
        dfli = casi
        dfl = cas
    }
}

if dfl != nil {
    selunlock(scases, lockorder)
    casi = dfli
    cas = dfl
    goto retc
}
```

## 阻塞

在读写单个chan的过程中, 读写goroutine都会被阻塞. 同样的, 读写多个chan, 当前goroutine 也要阻塞.

阻塞的方式是对所有的chan加锁, 即 sellock()函数. 它首先通过堆排序, 对chan根据地址大小进行排序, 然后顺序加锁. 采用堆排序的理由是, 堆排序时间复杂度为 $O(nlogn)$, 而空间复杂度为 $O(1)$.

接下来, 按照加锁顺序挨个改变对应的chan, 使其在等待队列中加入我们当前的 goroutine即可:

```go
switch cas.kind {
case caseRecv:
    c.recvq.enqueue(sg)

case caseSend:
    c.sendq.enqueue(sg)
}
```

## 为什么要按地址堆排序作为加锁顺序

前文中我们知道, select在有多个chan时, 会随机顺序遍历来寻找一个可用chan. 那么有没有办法让某个chan比别的有更大机会被读到呢?

方法当然是有的, select 允许多条 case 读写同一个ch. 

```go
select {
    case <-c0:
    // ... 0
    case <-c0:
    // ... 0
    case <-c1:
    // ... 1
    case <-c2:
    // ... 2
    case <-c3:
    // ... 3
    case <-c4:
    // ... 4
}
```

这样, 假如本来有5个不同的chan, 当5个ch均可读时, 每个chan被读取到的概率是 $\frac{1}{5}$ . 我们将c0多写了一遍, c0被读取到的概率变成了 $\frac{1}{3}$.

我们知道, 每个chan都有自己的读写goroutine队列, 读写时需要加锁. 而此时select读写同一个ch两遍, c0 的队列岂不是要被锁住两遍导致死锁, 或者队列中被加入两边同一个g?

在 sellock中对一系列select chan 进行加锁操作:

```go
func sellock(scases []scase, lockorder []uint16) {
	var c *hchan
	for _, o := range lockorder {
        c0 := scases[o].c
        // 比较是否c0与上一个相等
		if c0 != nil && c0 != c {
			c = c0
			lock(&c.lock)
		}
	}
}
```

由于我们已经根据地址排序了lockorder, 因此可以轻易的判断c0已经被锁过一次.
