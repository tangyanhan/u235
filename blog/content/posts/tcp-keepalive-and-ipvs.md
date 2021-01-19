---
title: "TCP保活机制闲扯"
date: 2021-01-19T22:20:28+08:00
tags: 网络
---

## TCP保活机制主要问题

在TCP连接不传输任何数据时，两端将无法得知另一端是否因为关机、掉电而消失，此外，为一个另一端已经不存在的TCP连接一直维持资源，也会浪费资源。在TCP标准之外，各个操作系统实现中都实现了TCP的保活机制(keepalive)。

保活机制有几个问题要解决：

1. 什么时候开始启动保活机制？

Linux 使用sysctl参数：`net.ipv4.tcp_keepalive_time=7200`，表示TCP连接在闲置7200秒后启动保活机制，开始发送保活报文

2. 保活报文间隔多久发送一次？

sysctl 参数: `net.ipv4.tcp_keepalive_intvl = 75`， 表示保活报文每隔75秒发送一次

3. 保活报文可能因为网络抖动而发送失败，失败多少次我们认为TCP连接已经需要放弃？

sysctl 参数: `net.ipv4.tcp_keepalive_probes = 9`，表示保活报文失败9次后，TCP连接被认为应当放弃，将会关闭

4. 开启保活机制的一方，可能会观察到哪些情形？

* 在保活过程中发生了数据传输，保活机制终止，等待下一次闲置触发

* 对方由于正在重启或中间的网络不可达，导致对报文没有任何响应，直到超时

* 对方重启成功，已经不记得连接的信息，因此返回`RST`报文，保活方收到`Connection Reset by Peer`

## 在Linux C/Golang中启用保活机制

保活机制不是TCP标准的一部分，通过系统调用创建的socket需要通过`setsockopt`明确启用`SO_KEEPALIVE`参数:

```c
#include <sys/socket.h>
/**level=SOL_SOCKET, option_name=SO_KEEPALIVE, option_value=1*/
int setsockopt(int socket, int level, int option_name,
const void *option_value, socklen_t option_len);
```

在Golang中，创建的TCP连接默认开启了KeepAlive（无论是Dial还是Accept）, 可以通过修改`Dialier.KeepAlive`时长修改保活报文间隔，当不设置它时，报文间隔是15s，而不是与Linux中的sysctl参数一致。

```go
// DialContext 主动建立连接时，默认会将TCP连接启用保活，并设置默认保活间隔
func (d *Dialer) DialContext(ctx context.Context, network, address string) (Conn, error) {
    // 很大一坨代码，略过
	if tc, ok := c.(*TCPConn); ok && d.KeepAlive >= 0 {
		setKeepAlive(tc.fd, true)
		ka := d.KeepAlive
		if d.KeepAlive == 0 {
			ka = defaultTCPKeepAlive
		}
		setKeepAlivePeriod(tc.fd, ka)
		testHookSetKeepAlive(ka)
    }
}
// Golang TCP accept中默认启用KeepAlive，设置了默认间隔
func (ln *TCPListener) accept() (*TCPConn, error) {
	fd, err := ln.fd.accept()
	if err != nil {
		return nil, err
	}
	tc := newTCPConn(fd)
	if ln.lc.KeepAlive >= 0 {
		setKeepAlive(fd, true)
		ka := ln.lc.KeepAlive
		if ln.lc.KeepAlive == 0 {
			ka = defaultTCPKeepAlive
		}
		setKeepAlivePeriod(fd, ka)
	}
	return tc, nil
}
```

## 保活机制相关的故障

从上面可知，在启用保活机制后，启用保活的一方总是能够得知对方是否已经“失踪”。但实际应用程序，往往在**应用层**忽视相关的错误处理，不进行错误检测和连接重建，导致已经断开的连接不被发现，等到发消息时才遇到报错一头雾水。

在生产中，有时会遇到一些特殊情形，例如“白天好好的，放一晚上，第二天早上就坏了”，往往就是因为TCP连接被静置两小时后（时长取决于Linux配置），应用层代码没能正确处理断开的TCP连接，导致下次发送数据时报错，或引发更多的错误。

在保活机制与IPVS结合之后，这类故障又有了新的呈现方式。

在Kubernetes中，当kube-proxy启用`mode: ipvs`时，使用ipvs实现cluster-ip到pod-ip的负载均衡，集群中到service的长连接，如GRPC，就会经由ipvs负载均衡到具体Pod上。

通过`ipvsadm -ln --timeout`命令，可以得到输出如下：
```
[root@cent1 ~]# ipvsadm -ln --timeout
Timeout (tcp tcpfin udp): 900 120 300
```

这表明，目前这台主机ipvs设定TCP超时时间为900秒，tcp fin为120秒（终止连接后等待的时间），udp超时时间为300秒。

那么，ipvs与保活机制会产生怎样的化合反应，出现什么样的大坑呢？

以我目前的设置来看，目前这台主机的TCP保活机制启动时间，是闲置7200秒以后：

```bash
[root@cent1 ~]# sysctl net.ipv4.tcp_keepalive_time
net.ipv4.tcp_keepalive_time = 7200
```

当我们的长连接通过IPVS连接到Pod之后，在闲置到900秒时，ipvs将会启动超时机制，将长连接断开。而保活机制要等到7200秒才能启动，因此长连接闲置到900秒时，总是会被不优雅的关闭。

这样在表面看来，就变成了“中午出去吃个饭，回来后应用开始出现故障”，因为900秒也就是15分钟而已。

简单的解决方式，就是将保活启动时间变短，使其小于ipvs的timeout。个人认为在应用层面写一些容错代码更好，因为总是期待运维人员去修改这些默认值，可能会让错误排查变得更加困难。

这个问题在这篇博客中有描述：[Kubernetes IPVS模式下服务间长连接通讯的优化，解决Connection reset by peer问题](https://blog.frognew.com/2018/12/kubernetes-ipvs-long-connection-optimize.html)
