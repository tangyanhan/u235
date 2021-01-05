---
title: "TCP重置报文段"
date: 2020-12-21T14:08:37+08:00
tags: "网络"
---

本节是作为《TCP/IP详解·卷1:协议》中TCP连接管理一章中“重置报文段”的读书笔记及一部分自己的扩展。

当TCP头部RST位字段置为1时，报文段就被称为“重置报文段”。在tcpdump监听流量时，会看到`R`字符标识。 这里介绍重置报文段的用途以及相关Linux配置。

# 重置报文段的用途

## 1. 拒绝不存在端口的连接请求

首先，启动tcpdump，监听本机lo接口:

![tcpdump](/images/post/tcp_syn/tcp_syn_conn_refused.png)

我们试着用curl访问本机上并未监听的端口8888，报告错误`Connection Refused`：

![connrefused](/images/post/tcp_syn/tcp_syn_conn_refused_2.png)

可以看到两条报文，第一条报文是curl发出的`SYN`，随后是本机做出的回应，其标志符为`R`，拒绝了我们的连接。

结论: 请求没有监听的端口，返回重置报文段，报错`Connection Refused`

## 2. 暴力终止连接

之前在TCP连接建立与终止中讲到，TCP正常的终止是使用`FIN`报文，这也被成为`有序释放`，它需要通信两端4次报文段交换才能结束。
但有时我们可以采用一种更暴力的手段来终止连接，比如服务端并不在意客户端是否收到了结束报文，因为它马上就要进行一次重启。通过`RST`报文终止连接，也被称为`终止释放`。终止发起方将会马上发出一个`RST`报文，并不再关注对方是否进行回复，而另一方收到时，会说明是收到了`RST`进行终止而非正常的`FIN`。

在Socket建立时可以通过`SO_LINGER`来设定这种行为。在Linux中，可以通过[setsockopt](https://linux.die.net/man/3/setsockopt)，将`SO_LINGER`设置为0，这样当调用`close`时，会直接发送重置报文段，不会等待缓存发送完成。在Golang中，对`TCPConn`对象调用`SetLinger(0)`即可。

我们写一段简单的Go服务端程序，将TCP连接通过`SetLinger(0)`将其设定为关闭时发送`RST`，然后尝试读写:

![conn_reset_by_peer](/images/post/tcp_syn/conn_reset_by_peer.png)
![conn_reset_by_peer_client](/images/post/tcp_syn/conn_reset_by_peer_client.png)

结论: 一方通过重置报文段直接终止了连接，另一端报错`Connection reset by peer`.

## 3. 半开连接

在通信过程中，通信一方关闭或终止连接，但不告知另一端，那么就会形成一个半开连接，例如一端电源被切断而非主动关机。在主机电源恢复后，由于已经不记得之前的TCP连接信息，对于另一端发来的TCP报文段，将会回复一个`RST`作为响应，另一端收到后将会主动关闭连接并反应情况。

## 4. 时间等待错误

正常关闭时，主动关闭方最后维持`TIME_WAIT`状态等待`2MSL`后彻底关闭，但此时如果收到一条重置报文段，就会破坏该状态，导致提前进入`CLOSED`状态。默认被动关闭方（服务器）此时连接进入`TIME_WAIT`状态，如果此时收到了之前的旧报文，由于连接已经关闭，服务器不记得信息，就会回复重置报文段，从而导致状态被破坏。许多系统实现时，连接处于`TIME_WAIT`状态时不会对重置报文段做出反应。

# Linux内核实现中的重置报文段

Linux中重置报文段的处理，在`net/ipv4/tcp_input.c:tcp_reset()`中:

```c
void tcp_reset(struct sock *sk)
{
	/* We want the right error as BSD sees it (and indeed as we do). */
	switch (sk->sk_state) {
	case TCP_SYN_SENT:
		sk->sk_err = ECONNREFUSED;
		break;
	case TCP_CLOSE_WAIT:
		sk->sk_err = EPIPE;
		break;
	case TCP_CLOSE:
		return;
	default:
		sk->sk_err = ECONNRESET;
	}
	/* This barrier is coupled with smp_rmb() in tcp_poll() */
	smp_wmb();

	tcp_write_queue_purge(sk);
	tcp_done(sk);
```

1. 当主动建立请求时，得到`RST`，则报错 `ECONNREFUSED： Connection refused`
2. 当服务端连接处于`CLOSE_WAIT`状态，收到`RST`,则报错 `EPIPE: Broken pipe`
3. 连接已经关闭，处于`CLOSED`状态，忽略
4. 其他情况，一律报错 `ECONNRESET: Connection reset by peer`

