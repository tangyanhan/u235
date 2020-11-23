---
title: "eBPF入门-2: bpftrace"
date: 2020-07-04T20:24:33+08:00
tags: "eBPF"
---

本文由我翻译自[The bpftrace One-Liner Tutorial](https://github.com/iovisor/bpftrace/blob/master/docs/tutorial_one_liners.md). 

# 安装

```bash
sudo apt install -y bpftrace
```

# 1. 列出可用的追踪

```bash
sudo bpftrace -l 'tracepoint:syscalls:sys_enter_*'
```

bpftrace能够跟踪很多系统调用，这里是可以用 * 和 ? 查询的。

# 2. Hello World

```bash
ethan@ethan-kali:~$ sudo bpftrace -e 'BEGIN {printf("hello world\n"); }'
Attaching 1 probe...
hello world
^C
```

* BEGIN 就像awk的BEGIN一样， 可以用于设置变量， 打印事件头等
* 跟踪过程可以关联一些动作，它们定义在花括号里。

# 3. 跟踪打开的文件

```bash
ethan@ethan-kali:~$ sudo bpftrace -e 'tracepoint:syscalls:sys_enter_openat { printf("%s %s\n", comm, str(args->filename)); }'
Attaching 1 probe...
code /proc/6972/cmdline
ksysguardd /proc/stat
ksysguardd /proc/vmstat
ksysguardd /proc/meminfo
```

* 开始是 ```tracepoint:syscalls:sys_enter_openat```， 这是我们要监听的系统调用。在现代的Linux系统中(glibc>=2.26)，open总是会调用openat.

* comm 是一个内置的变量名，也就是当前进程的名字。类似的还有 pid/tid.

* args是一个指针，指向一个包含了所有跟踪点参数的结构体。 我们可以通过这样的命令找出结构体中的所有成员： ```bpftrace -vl tracepoint:syscalls:sys_enter_openat```

* args->filename 解引用 args 结构体并获得其中的 filename 成员。

* str() 将指针指向的字符串转换成字符串（因为上面的filename也是个字符串指针)

查看 ```sys_enter_openat``` 的args包含了哪些成员：

```bash
ethan@ethan-kali:~$ sudo bpftrace -vl tracepoint:syscalls:sys_enter_openat 
tracepoint:syscalls:sys_enter_openat
    int __syscall_nr;
    int dfd;
    const char * filename;
    int flags;
    umode_t mode;
```

# 4. 统计进程系统调用数

```bash
ethan@ethan-kali:~$ sudo bpftrace -e 'tracepoint:raw_syscalls:sys_enter { @[comm] = count(); }'
Attaching 1 probe...
^C

@[haveged]: 1
@[kglobalaccel5]: 2
@[org_kde_powerde]: 2
@[DiscoverNotifie]: 2
@[kdeconnectd]: 2
@[EventHandler]: 2
```

* @作用是声明了一个map。 @后面可以跟一个名字，这样可以把不同的map区分开来（因为实际是可以使用多个map的）

* [] 设定了map关联的key，这里使用了comm，即进程名作为key。

* count(): 这是一个map 函数，决定了如何统计map。 count() 本身返回系统调用的次数，这里由于key是进程名， 因此返回根据进程名统计的系统调用次数

当bpftrace结束（比如按下Ctrl+C)，就会自动打印出所有的map。

# 5. 打印 read() 读取字节数的分布

```bash
ethan@ethan-kali:~$ sudo bpftrace -e 'tracepoint:syscalls:sys_exit_read /pid == 6615/ { @bytes = hist(args->ret); }'
Attaching 1 probe...
^C

@bytes: 
(..., 0)               9 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
[0]                    1 |@@@@@                                               |
[1]                    0 |                                                    |
[2, 4)                 0 |                                                    |
[4, 8)                 0 |                                                    |
[8, 16)                1 |@@@@@                                               |
[16, 32)               0 |                                                    |
[32, 64)               2 |@@@@@@@@@@@                                         |
[64, 128)              0 |                                                    |
[128, 256)             2 |@@@@@@@@@@@                                         |
[256, 512)             1 |@@@@@                                               |
[512, 1K)              3 |@@@@@@@@@@@@@@@@@                                   |
[1K, 2K)               2 |@@@@@@@@@@@                                         |
```

以上命令，统计了PID为6615的程序在此期间调用 sys_read() 的返回值， 并打印出统计柱形图。

* /.../ : 这是过滤器，只有符合过滤器才能继续下面的动作(action)。这里的设置是只监听pid==6655的进程相关活动。类似的， "&&" 和 "||" 操作符也是支持的。

* ret: 这是函数的返回值。对于```sys_read```来说， 是-1(错误)或者成功读取到的字节数

* @: 这是一个类似上面一节的map，但没有任何的键值([])，bytes定义了map的名字，使之能够与其它map打印时区别开来

* hist(): 这是一个map函数， 以2的幂次作为区间作出分布图。打印出的第一列是统计区间(即统计图中的X轴上区间），下一列是出现次数， 下一列的@@@@则是统计次数的图形化表示。

* 其它的map函数还有 lhist(linear hist), count(), sum(), avg(), min(), 以及 max()

# 6. 动态跟踪内核中 read() 的字节数

```bash
ethan@ethan-kali:~$ sudo bpftrace -e 'kretprobe:vfs_read { @bytes = lhist(retval, 0, 2000, 200); }'
Attaching 1 probe...
^C

@bytes: 
(..., 0)              44 |@                                                   |
[0, 200)            1415 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
[200, 400)             0 |                                                    |
[400, 600)             0 |                                                    |
[600, 800)             0 |                                                    |
[800, 1000)            0 |                                                    |
[1000, 1200)          18 |                                                    |
[1200, 1400)           0 |                                                    |
[1400, 1600)           0 |                                                    |
[1600, 1800)           0 |                                                    |
[1800, 2000)           0 |                                                    |
[2000, ...)            3 |                                                    |
```

ET: 这里使用了上一节提到的lhist, 即线性统计，与hist的2幂次统计区间不同， lhist 可以指定统计的起始及终止区间， 以及步长：

```
lhist(value, min, max, step)
```

* 开头的 ```kretprobe:vfs_read``` ， 首先 kretprobe 表示用于追踪内核函数的返回值， 随后 vfs_read 表示追踪的函数是 vfs_read. 由于内核版本不同， 其中各种函数的名字，参数，返回值都可能发生变化， 因此会发生改变。除了 kretprobe, 还有其它的如追踪函数开始执行。 这些追踪很强大， 但需要对内核源码有足够了解。

# 7. 统计 read() 耗时

```
ethan@ethan-kali:~$ sudo bpftrace -e 'kprobe:vfs_read { @start[tid] = nsecs; } kretprobe:vfs_read /@start[tid]/ { @ns[comm] = hist(nsecs - @start[tid]); delete(@start[tid]); }'
Attaching 2 probes...
^C

@ns[kactivitymanage]: 
[4K, 8K)               1 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|

@ns[xembedsniproxy]: 
[4K, 8K)               1 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|

@ns[systemd-journal]: 
[4K, 8K)               1 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|

@ns[code]: 
[4K, 8K)               2 |@@@@@@@@@@@@@                                       |
[8K, 16K)              3 |@@@@@@@@@@@@@@@@@@@                                 |
[16K, 32K)             8 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|                       |

@start[23580]: 91727748589219
@start[1385]: 91727748883680
@start[1212]: 91727749317256
@start[23582]: 91727749318410
@start[1235]: 91727749799433
@start[30483]: 91727781856565
```

* @start[tid]: 使用Thread ID作为键值。 可能有很多正在进行中的 read()， 我们想记下每个开始的时间戳。 由于内核中线程每次只能进行一个系统调用， 因此这里使用了tid作为键值。

* nsecs: 自从启动以来的纳秒，可以用于时间事件中的高精度时间戳。

* /@start[tid]/: 这个过滤器确保只处理我们记录过开始时间的thread。 这是为了避免我们启动跟踪时， 去跟踪某个已经开始的read()，这样会由于我们没有记录过开始时间而出错。

* delete(@start[tid]): 释放变量。

# 8. 统计进程级事件

```
ethan@ethan-kali:~$ sudo bpftrace -e 'tracepoint:sched:sched* { @[probe] = count(); } interval:s:5 { exit(); }'
Attaching 25 probes...


@[tracepoint:sched:sched_process_wait]: 11
@[tracepoint:sched:sched_migrate_task]: 1575
@[tracepoint:sched:sched_wake_idle_without_ipi]: 29166
@[tracepoint:sched:sched_wakeup]: 37333
@[tracepoint:sched:sched_stat_runtime]: 37754
@[tracepoint:sched:sched_waking]: 37911
@[tracepoint:sched:sched_switch]: 73204
```

* sched: sched 跟踪类别下有一系列高级的调度其及进程事件， 例如 fork, exec, context switch(进程切换)
* probe: 跟踪的全名（完整的系统调用函数名）
* interval:s:5 : 作用是每5秒进行一次统计，可以用作统计间隔或超时。
* exit(): 使bpftrace在一次统计后退出。

# 9. 统计CPU上的内核栈性能

```
ethan@ethan-kali:~$ sudo bpftrace -e 'profile:hz:99 { @[kstack] = count(); }'
Attaching 1 probe...
^C

@[
    pcibios_pm_ops+409080112
]: 1
@[
    run_timer_softirq+64
    __softirqentry_text_start+230
    irq_exit+166
    smp_apic_timer_interrupt+118
    apic_timer_interrupt+15
]: 1
@[
    find_busiest_group+227
    load_balance+370
    rebalance_domains+684
    __softirqentry_text_start+230
    irq_exit+166
    smp_apic_timer_interrupt+118
    apic_timer_interrupt+15
    cpuidle_enter_state+201
    cpuidle_enter+41
    do_idle+484
    cpu_startup_entry+25
    start_secondary+351
    secondary_startup_64+164
]: 1
```

以99Hz频率分析内核栈性能， 打印出频率统计。

* profile:hz:99 : 在所有CPU上以99Hz进行分析。 为什么不是100或1000？ 因为我们希望让频率够快，能够同时捕捉到运行中大和小的运行情况，但不能过于频繁影响性能。 100Hz是够的，但我们不希望100, 因为取样可能会发生在其它计时活动中， 因此99.

* kstack: 返回内核调用栈。这个被用作map的键值，所以可以统计其频率。这个输出适合图形化做成火焰图。 追踪用户空间的栈， 可以用 ```ustack```。

# 10. 跟踪调度器

```
ethan@ethan-kali:~$ sudo bpftrace -e 'tracepoint:sched:sched_switch { @[kstack] = count(); }'
[sudo] password for ethan: 
Attaching 1 probe...
^C

@[
    __sched_text_start+866
    __sched_text_start+866
    schedule+74
    schedule_timeout+350
    intel_pipe_update_start+337
    intel_update_crtc+161
    skl_commit_modeset_enables+437
    intel_atomic_commit_tail+779
    process_one_work+436
    worker_thread+80
    kthread+249
    ret_from_fork+53
]: 1
@[
    __sched_text_start+866
    __sched_text_start+866
    schedule+74
    schedule_timeout+350
    rtR0SemEventLnxWait.isra.0+835
    supR0SemEventWaitEx+123
    supdrvIOCtl+10757
    VBoxDrvLinuxIOCtl_6_1_10+362
    ksys_ioctl+135
    __x64_sys_ioctl+22
    do_syscall_64+82
    entry_SYSCALL_64_after_hwframe+68
]: 1
```

统计导致上下文切换事件的内核栈。

* sched: sched类别提供了不同CPU调度器事件的跟踪点： sched_switch, sched_wakeup, sched_migrate_task, 等等。
* sched_switch: 这个跟踪点在线程离开CPU时触发，其原因可能是一个等待事件， 例如等待I/O, 时钟， 分页中断或交换， 或者是一个锁。
* kstack: 将内核栈用作map的键值。
* sched_switch 是在进程上下文切换时触发的，因此栈此时描述的是要离开的线程。当你在使用其它类型的追踪时，要注意上下文，可能追踪中的 comm, pid, kstack,等等， 可能并不是追踪的目标。

# 11. 阻塞性 I/O 追踪


```
ethan@ethan-kali:~$ sudo bpftrace -e 'tracepoint:block:block_rq_issue { @ = hist(args->bytes); }'
Attaching 1 probe...
^C

@: 
[0]                   27 |@@@@@@                                              |
[1]                    0 |                                                    |
[2, 4)                 0 |                                                    |
[4, 8)                 0 |                                                    |
[8, 16)                0 |                                                    |
[16, 32)               0 |                                                    |
[32, 64)               0 |                                                    |
[64, 128)              0 |                                                    |
[128, 256)             0 |                                                    |
[256, 512)             0 |                                                    |
[512, 1K)            117 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@                        |
[1K, 2K)              25 |@@@@@@                                              |
[2K, 4K)               7 |@                                                   |
[4K, 8K)             215 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
[8K, 16K)             51 |@@@@@@@@@@@@                                        |
[16K, 32K)            32 |@@@@@@@                                             |
[32K, 64K)             2 |                                                    |
[64K, 128K)            3 |                                                    |
[128K, 256K)          21 |@@@@@                                               |
```

追踪阻塞性I/O请求，统计字节数量。

* tracepoint:block: block分类，主要用于跟踪各种阻塞性I/O(存储)事件。
* block_rq_issue: 该事件当阻塞性I/O请求指派(issued)到设备时触发。
* args->bytes: 从 block_rq_issue 参数中取出，是请求的字节数。

这个探查的上下文很重要： 当阻塞性I/O请求指派到设备时， 发生该事件。 这一般发生在进程上下文中， 这时comm能够提供给你进程的名字。 但也可能发生在内核上下文中， 例如 readahead， 这时 pid/comm 都不会提供给你符合预期的应用程序。

# 12. 内核结构追踪

创建并编辑文件```path.bt```: 

```c
#include <linux/path.h>
#include <linux/dcache.h>

kprobe:vfs_open
{
        printf("open path: %s\n", str(((struct path *)arg0)->dentry->d_name.name));
}
```

随后运行命令：

```bash
ethan@ethan-kali:~$ sudo bpftrace path.bt 
Attaching 1 probe...
open path: cmdline
open path: cmdline
open path: status
```

这里动态追踪了 vfs_open() 函数，该函数地一个参数为 (struct path*)。

* kprobe: 内核动态追踪调查类型， 分类下包含了各种内核函数入口事件（对比kretprobe)
* arg0: 一个内置变量， 指向地一个调查的参数， 具体如何取决于参数原型。 其他参数可以通过 arg1, ..., argN 访问。
* ```((struct path *)arg0)->dentry->d_name.name```: 将arg0 转换为 ```struct path *```，然后将其解引用, 等等。
* #include: 用于包含结构体 path/dentry的定义，好让我们能够使用它们。

内核结构体的支持与bcc相同， 都是使用linux 的内核头文件。 这意味着许多结构提都可用，不过不是全部， 有时可能需要手动的include一个结构体。 对此， 可以参考 [dcsnoop tool](https://github.com/iovisor/bpftrace/blob/master/tools/dcsnoop.bt)， 包含了一部分手动引入的结构提 nameidata， 因为在头文件中找不到它。 如果内核有BTF (BPF Type Format)数据， 所有内核结构都始终可用。

现在你已经对bpftrace有了较多的了解， 也可以开始使用和写自己的一行命令了。 要了解更多内容， 参考 [Reference Guide](https://github.com/iovisor/bpftrace/blob/master/docs/reference_guide.md)