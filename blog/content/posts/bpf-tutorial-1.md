---
title: "eBPF入门-1: 概要，bcc工具包"
date: 2020-06-29T11:43:33+08:00
tags: "eBPF"
---

# 什么是BPF，eBPF

BPF全程是 Berkeley Packet Filters， 由McCanne于1992年提出，并加入到Unix中。 从名称可以看出， 这就是个包过滤器。它是 libpcap/tcpdump/wireshark这些 Linux 嗅探器包的基石。

BPF在Linux内核中由两个部分组成： Network Tap 和 过滤器。Network Tap就像一个水龙头分流阀一样， 从物理接口传来的数据包，都会同时被Network Tap和网络栈(Network Stack)各自处理，符合过滤器条件的数据包， 会被BPF复制到缓冲区，提供给对应的用户空间程序， 例如tcpdump。

过滤器部分本质上是一个虚拟机(Pseudo Machine), 它规定了一种不依赖具体协议的过滤器语言，能够将用户传来的代码编译成有限步骤的高效过滤器。

![BPF Overview](/images/post/bpf-overview.png)

过滤器允许用户通过系统调用，将一段编译好的代码嵌入到内核态的BPF过滤器中，成为一个新的过滤器，这样过滤器就好像TCP网络栈一样，可以让应用程序在用户态直接源源不断的获得到网络数据包。

从这个架构可以看出， BPF可以用于流量监听， 但不能影响网络栈的处理过程。 不能影响正常的应用程序处理。

eBPF，即Extened BPF， 原理与BPF基本一致，但范围大大扩展了，不单独是网络流量， 系统调用等信息也可以流入BPF的虚拟机中。

# eBPF可以做什么？

* 经典的BPF部分，可以做网络流量分析
* 经典的BPF部分，必要时也可以做流量复制转发
* 新增的eBPF，可以做系统行为分析和系统性能/应用程序性能分析等

# 安装bcc工具包

```bash
sudo apt install -y bpfcc-tools
```

实验一下是否安装好了：

```bash
sudo /usr/sbin/opensnoop-pbfcc
```

opensnoop 是bcc工具包的一部分，作用是监听指定或所有进程的open()系统调用，我们可以获得进程的文件访问行为，用来进行进程行为分析等。

# bcc工具包使用

以下部分由我翻译自 [bcc Tutorial](https://github.com/iovisor/bcc/blob/master/docs/tutorial.md)

## opensnoop

opensnoop 监听所有open()系统调用行为。可以从中分析出进程的数据文件，日志文件等。也可以判断出进程的违规行为，或者某些使用不当导致的性能问题，例如频繁访问一个不存在的文件。

## execsnoop

execsnoop 监听所有exec()系统调用行为。需要注意的是它监听 exec()而不是fork()，因此并不能监听所有新起的进程。

```bash
ethan@ethan-kali:~$ sudo execsnoop-bpfcc 
PCOMM            PID    PPID   RET ARGS
chrome-sandbox   11239  2300     0 /opt/google/chrome/chrome-sandbox --adjust-oom-score 11238 300
```

## ext4slower(或者 brtfsslower, xfsslower, zfsslower)

用于监听对应文件系统的访问行为，找出进程缓慢的文件系统访问操作。

```bash
ethan@ethan-kali:~$ sudo xfsslower-bpfcc 
Tracing XFS operations slower than 10 ms
TIME     COMM           PID    T BYTES   OFF_KB   LAT(ms) FILENAME
15:00:40 b'code'        6480   S 0       0          49.63 b'19b178a2d49de9d093180fc06588a7e5'
15:00:48 b'code'        6480   S 0       0          10.17 b'state.vscdb'
```

## biosnoop

biosnoop主要用于打印出每次磁盘IO及延迟。

```bash
ethan@ethan-kali:~$ sudo biosnoop-bpfcc 
TIME(s)     COMM           PID    DISK    T SECTOR     BYTES  LAT(ms)
0.000000    kworker/4:2    7497           R 18446744073709551615 8         0.53
0.003048    AioMgr0-N      2618   nvme0n1 R 395545048  4096      4.95
0.003368    AioMgr0-N      2618   nvme0n1 R 395547840  8192      0.11
0.003398    AioMgr0-N      2618   nvme0n1 R 394532376  4096      0.13
0
```

## cachestat

cachestat 用于分析每秒（或指定间隔）文件系统缓存的命中率。可以通过分析命中率和失误率，分析文件系统缓存对性能的影响。

```bash
ethan@ethan-kali:~$ sudo cachestat-bpfcc 
    HITS   MISSES  DIRTIES HITRATIO   BUFFERS_MB  CACHED_MB
      26        0        0  100.00%            3       3512
      20        0        0  100.00%            3       3508
      20        0        0  100.00%            3       3508
      20        0        0  100.00%            3       3508
```

## tcpconnect

通过跟踪connect()系统调用，每次有活跃TCP连接建立时打印， 可以用于分析出程序错误配置导致TCP连接低效使用，或者是入侵等问题。

```bash
ethan@ethan-kali:~$ sudo tcpconnect-bpfcc 
Tracing connect ... Hit Ctrl-C to end
PID    COMM         IP SADDR            DADDR            DPORT
2343   Chrome_Child 4  127.0.0.1        127.0.0.1        1080
2343   Chrome_Child 4  127.0.0.1        127.0.0.1        1080
10707  shadowsocks- 4  10.8.1.36        47.103.69.196    59012
10707  shadowsocks- 4  10.8.1.36        47.103.69.196    59012
```

## tcpaccept

同上， 只是追踪外界主动连接过来的请求。

## tcpretrans

```bash
ethan@ethan-kali:~$ sudo tcpretrans-bpfcc 
Tracing retransmits ... Hit Ctrl-C to end
TIME     PID    IP LADDR:LPORT          T> RADDR:RPORT          STATE
15:13:40 0      4  10.8.1.36:59540      R> 180.149.145.248:443  LAST_ACK
15:13:41 0      4  10.8.1.36:59540      R> 180.149.145.248:443  LAST_ACK
1
```

tcpretrans 打印出每个重传的TCP包。TCP重传会导致延迟和吞吐量问题。 

ESTABLISHED 重传， 可能由于网络问题导致。SYN_SENT 可能是由于CPU占用过高或内核丢包过多导致的。

TODO: 搞清楚上述原理。

## runqlat

runqlat 用于跟踪一段时间内进程/线程的调度延迟，开始一段时间后， 按下Ctrl+C, 可以打印出这段时间的CPU调度延迟分布图。

```bash
ethan@ethan-kali:~$ sudo runqlat-bpfcc 
Tracing run queue latency... Hit Ctrl-C to end.
^C
     usecs               : count     distribution
         0 -> 1          : 624175   |**********************************      |
         2 -> 3          : 514640   |****************************            |
         4 -> 7          : 541942   |*****************************           |
         8 -> 15         : 285838   |***************                         |
        16 -> 31         : 185164   |**********                              |
        32 -> 63         : 723860   |****************************************|
        64 -> 127        : 56832    |***                                     |
       128 -> 255        : 18260    |*                                       |
       256 -> 511        : 9960     |                                        |
       512 -> 1023       : 5129     |                                        |
      1024 -> 2047       : 1916     |                                        |
      2048 -> 4095       : 634      |                                        |
      4096 -> 8191       : 137      |                                        |
      8192 -> 16383      : 26       |                                        |
```

解读：

usecs 表明这一列时间单位是us. count 即发生对应调度延迟的次数， distribution 是将count进行了更直观的图形化。

由上图可知， 过去这段时间， 有18260次线程调度延迟在 128-255 us 之间。

我们可以由此分析出花在CPU调度上的延迟， 从而决定水平扩展或物理扩展来提高性能。

## profile

```bash
ethan@ethan-kali:~$ sudo profile-bpfcc
Sampling at 49 Hertz of all threads by user + kernel stack... Hit Ctrl-C to end.

    b'pcibios_pm_ops'
    b'pcibios_pm_ops'
    ioctl
    [unknown]
    [unknown]
    [unknown]
    [unknown]
    [unknown]
    [unknown]
    start_thread
    -                EMT-0 (12255)
        1

    fbGetScreenPrivateKey
    [unknown]
```

profile 可以打印出进程的栈和系统调用，以及他们的发生次数， 可以用于分析进程行为。


# 参考文档

1. BPF论文： [The BSD Packet Filter:A New Architecture for User-level Packet Capture](https://www.tcpdump.org/papers/bpf-usenix93.pdf)