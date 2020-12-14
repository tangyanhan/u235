---
title: "TCP协议基础:连接建立，关闭与状态转移"
date: 2020-12-03T18:40:15+08:00
tags: "网络"
---

本篇是对TCP/IP详解的相关读书笔记。

# TCP连接的建立与终止

**一般情况下**，C/S通讯模型下的TCP连接的建立与关闭可以概括为“三次握手，四次关闭”。

## 三次握手

1. Client发起SYN，并带上自己随机的初始化ISN(c)
2. Server收到后, 回复SYN+ACK，通过设定Seq=ISN(c)+1，表示自己已经正确收到了该SYN，并带上自己的ISN(s)
3. Client回复ACK, 设定Seq=ISN(c)+1,表示这是第一个成功发送的包，ACK=ISN(s)+1，表示自己正确收到了第二部的SYN

ISN是一个随机的16位二进制数字，通信两侧随机选择一个数字作为自己的初始化ISN.

假设两侧固定选择0，那么会发生什么呢？ 由于TCP连接的建立在网络上没有加密和其它验证机制，骚扰者可以通过不断伪造包来打断两者之间的握手过程。

## 四次断开

由于TCP协议是一个全双工协议，通讯双方都可以主动向对方发送数据，因此关闭时需要明确的关闭双方通道，共需`(FIN+ACK)X2`，即四次断开。

1. 客户端发送`FIN+ACK, Seq=K, ACK=L`, 这里ACK时对上一条报文的回复，Seq是客户端的当前计数
2. 服务端收到后，回复`ACK,Seq=L, ACK=K+1`，`ACK`表明自己正确收到了上述信息，并回应。此时客户端发送通道已经关闭，服务端仍可继续向客户端发送剩余的缓冲信息
3. 服务端发送`FIN+ACK, Seq=L, ACK=K+1`
4. 客户端回复`ACK,Seq=K, ACK=L+1`,表明自己正确收到了关闭信息，服务端收到后将释放相关资源

## 三次连接/四次断开示意图

![TCP连接的建立与终止](/images/post/tcp_states.png)

## 同时打开

同时打开不是指A通过客户端请求B的端口7777，而B同时通过客户端请求A的端口8888。 因为默认情况下，通过客户端connect连接到主机特定端口，客户端会随机选择一个端口。这时，两者同时建立了两个不同的TCP连接。

同时打开，是指A通过8888端口连接B的7777端口同时，B也通过7777端口连接A。由于两者选择的都是相同的端口组合，因此建立的是同一个TCP连接。由于两边都是主动打开者，根据三次握手的设定，主动打开连接时，对方必须回应`SYN+ACK`才能让连接正常建立，因此同时打开，会多一次`SYN+ACK`。

## 同时关闭

在只有一方主动关闭TCP连接的情况下，TCP连接是“4次关闭”，即主动关闭一方发出`FIN`，被动关闭方发出`FIN+ACK`表示收到了对方的`FIN`，然后也发出自己的`FIN`表示关闭自己这边的半连接，主动关闭方回复`FIN+ACK`后，连接终止。

在双方同时关闭的情况下，两边都是主动关闭方，几乎同时发出了自己的`FIN`表示要关闭自己的半连接，在收到对方的`FIN`时，他们各自的`FIN`已经发出，因此只需要对对方的`FIN`回复`FIN+ACK`就完成了两个半连接的关闭。同时关闭依旧是4次报文交换，但报文段的顺序是交叉的。

# TCP的状态转移

综上所述，TCP连接的建立和终止实际包括了主动+被动，以及双方都是主动的同时打开，同时关闭等。TCP中这些报文会引起连接进入不同的状态，它是一个有限状态机。
在[RFC 793](https://tools.ietf.org/html/rfc793)中描述了其状态转移，这里照搬TCP/IP详解中的状态转移图:

![TCP状态转移图](/images/post/tcp_state_machine.png)

## TIME_WAIT状态

`TIME_WAIT`又成为`2MSL`等待状态，因为在该状态时，TCP连接总是会等待两倍于最大段生存周期(Max Segment Lifetime, MSL)的时间。一个MSL是指一个报文段在网络中被丢弃之前，最大存活的周期。在Linux中，`2MSL`这个值可以通过查看sysctl参数`net.ipv4.tcp_fin_timeout`得到：

```bash
sysctl net.ipv4.tcp_fin_timeout
# 或者
cat /proc/sys/net/ipv4/tcp_fin_timeout
```

对于IPv6而言，Linux同样使用该参数来控制这个值，没有额外的键值。

### TIME_WAIT存在的意义

1. 在主动关闭连接后，为了防止我们最后的ACK丢失，等待一段时间

2. 在等待期间，双方将该连接（客户端IP地址，客户端端口号，服务器IP地址，服务端端口号）标记为不可用，避免后续新连接与前面的连接混淆。

`TIME_WAIT`状态下，如果复用TCP连接可能出现混淆，因为我们难以单凭ISN区分出不同的连接，ISN可能出现环绕，也可能出现两个连接重叠，因此在[RFC6191](https://tools.ietf.org/html/rfc6191)中定义了通过定义时间戳来更精确的区分不同连接的TCP报文段，从而复用`TIME_WAIT`状态的TCP连接。

因此，在Linux下，复用`TIME_WAIT`状态连接，需要同时开启`net.ipv4.tcp_tw_reuse`及`net.ipv4.tcp_timestamps`：
```
sysctl -w net.ipv4.tcp_tw_reuse=1
sysctl -w net.ipv4.tcp_timestamps=1
```

## FIN_WAIT_2状态

从状态转移图可知，`FIN_WAIT_2`出现在主动关闭一方，当主动关闭方发送了自己的`FIN`并收到`FIN+ACK`后，只要在等待对方的`FIN`，就可以回复`ACK`并进入`TIME_WAIT`状态，这就是`FIN_WAIT_2`，即等待4次关闭中的第二次`FIN`，而这时被动关闭方将停留在`CLOSE_WAIT`状态，直到发出自己的`FIN`。为了防止双方因为没有应答而永远卡在半关闭状态，Linux会在连接超时且没有收到回复时，将连接直接转入`CLOSED`状态，这个值是`net.ipv4.tcp_fin_timeout`。

# TCP连接管理相关攻击

## SYN泛洪

即若干主机大量向主机发送`SYN`而不进一步建立连接，导致大量浪费主机资源而导致拒绝服务。这里可以使用一种叫做`SYN Cookies`的技术，通过对SYN建档和回复巧妙设计的`SYN+ACK`序列号，检查是否客户端进行了正确的回复，然后才为连接分配相应的资源。

Linux下开启该功能参数为:
```
sysctl -w net.ipv4.tcp_syncookies=1
```

# TCP在Linux系统中的状态表示

如果使用如lsof, netstat这些工具，会发现有一列状态，这些状态码含义在内核源码 `net/ipv4/tcp.c` 中有描述:

```c
/*
 *	TCP_SYN_SENT		发出SYN请求建立连接
    sent a connection request, waiting for ack

 *	TCP_SYN_RECV		收到SYN连接请求，发送ACK，并等待三次握手的最后一个ACK
    received a connection request, sent ack,
 *				waiting for final ack in three-way handshake.
 *
 *	TCP_ESTABLISHED		连接建立/connection established
 *
 *	TCP_FIN_WAIT1		我方关闭，等待剩余缓存数据发送完成
    our side has shutdown, waiting to complete
 *				transmission of remaining buffered data
 *
 *	TCP_FIN_WAIT2		所有缓存数据发送完成，等待远程关闭
    all buffered data sent, waiting for remote
 *				to shutdown
 *
 *	TCP_CLOSING		两侧都已经关闭，但我们还有必须要发送完的数据
    both sides have shutdown but we still have
 *				data we have to finish sending
 *
 *	TCP_TIME_WAIT	进入关闭状态钱等待最后一个ACK，只能从FIN_WAIT2或CLOSING进入该状态
  	timeout to catch resent junk before entering
 *				closed, can only be entered from FIN_WAIT2
 *				or CLOSING.  Required because the other end
 *				may not have gotten our last ACK causing it
 *				to retransmit the data packet (which we ignore)
 *
 *	TCP_CLOSE_WAIT	远端已关闭，等待我们把剩下的数据写入并关闭，我们必须调用close()来进入到LAST_ACK状态
 	remote side has shutdown and is waiting for
 *				us to finish writing our data and to shutdown
 *				(we have to close() to move on to LAST_ACK)
 *
 *	TCP_LAST_ACK 对方已经关闭后，我方也已经关闭。或许我们的缓存里还有需要发送的数据.
 	out side has shutdown after remote has
 *				shutdown.  There may still be data in our
 *				buffer that we have to finish sending
*/
```
