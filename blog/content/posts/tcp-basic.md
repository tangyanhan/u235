---
title: "TCP协议基础:连接建立，关闭与状态转移"
date: 2020-12-03T18:40:15+08:00
tags: "网络"
---

# TCP连接的建立与关闭

**一般情况下**，C/S通讯模型下的TCP连接的建立与关闭可以概括为“三次握手，四次关闭”。

三次握手：

1. Client发起SYN，并带上自己随机的初始化ISN(c)
2. Server收到后, 回复SYN+ACK，通过设定Seq=ISN(c)+1，表示自己已经正确收到了该SYN，并带上自己的ISN(s)
3. Client回复ACK, 设定Seq=ISN(c)+1,表示这是第一个成功发送的包，ACK=ISN(s)+1，表示自己正确收到了第二部的SYN

ISN是一个随机的16位二进制数字，通信两侧随机选择一个数字作为自己的初始化ISN.

假设两侧固定选择0，那么会发生什么呢？ 由于TCP连接的建立在网络上没有加密和其它验证机制，骚扰者可以通过不断伪造包来打断两者之间的握手过程。

四次断开：

由于TCP协议是一个全双工协议，通讯双方都可以主动向对方发送数据，因此关闭时需要明确的关闭双方通道，共需`(FIN+ACK)X2`，即四次断开。

1. 客户端发送`FIN+ACK, Seq=K, ACK=L`, 这里ACK时对上一条报文的回复，Seq是客户端的当前计数
2. 服务端收到后，回复`ACK,Seq=L, ACK=K+1`，`ACK`表明自己正确收到了上述信息，并回应。此时客户端发送通道已经关闭，服务端仍可继续向客户端发送剩余的缓冲信息
3. 服务端发送`FIN+ACK, Seq=L, ACK=K+1`
4. 客户端回复`ACK,Seq=K, ACK=L+1`,表明自己正确收到了关闭信息，服务端收到后将释放相关资源

!!! 状态转移图存疑，需要一定修正

![状态转移示意图](/static/post/tcp_states.png)

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
