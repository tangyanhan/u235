---
title: "Go HTTP连接复用"
date: 2020-12-01T23:21:33+08:00
---

HTTP 1.1中给出了连接复用的方法，即通过设定Header为`Connection: keep-alive`，服务端如果支持此选项，那么会在返回中同样设置该Header，请求结束后不会立即关闭连接。

HTTP的连接复用与TCP的reuse是两回事，两个使用不同的机制实现。

这里描述的“连接”，包括TCP连接，也包含其上的TLS连接，因此HTTP的Keep Alive实现的连接复用，省去了TCP连接建立以及TLS连接建立的过程。

# 服务端

```go
// If this was an HTTP/1.0 request with keep-alive and we sent a
// Content-Length back, we can make this a keep-alive response ...
if w.wants10KeepAlive && keepAlivesEnabled {
	sentLength := header.get("Content-Length") != ""
	if sentLength && header.get("Connection") == "keep-alive" {
		w.closeAfterReply = false
	}
}
```

# 客户端

Golang的HTTP Client通过`net/http/trasnport.go`中的`Transport`对象实现底层TLS/TCP连接的封装。在`Transport`中，主要有以下几个参数：

* `DisableKeepAlives`: 默认为false，如果设为`true`,那么所有连接复用的优化选项都无效
* `MaxIdleConns`: 最大空闲连接数，该Transport可以维护最大这么多的空闲连接，用于连接复用, 为0时表示无限制
* `MaxIdleConnsPerHost`: 连接到每个host的最大闲置连接数，如果为0，就会使用`DefaultMaxIdleConnsPerHost`，这个值在go1.15是2
* `MaxConnsPerHost`: 连接单个host最大连接数，如果超了，那么超出的连接要等待
* `IdleConnTimeout`: 如果一个连接一段时间没有用，那么会由客户端主动关闭，为0时表示没有限制

在Dial建立连接后，就会开始进行读循环和写循环。在读循环中，能够获得HTTP Response，其中包括Header以及Body。当Body被读至末尾EOF，或者被手动关闭时，这个connection就被视为idle，可以回收用于其它请求了。

在knative的[kafkasource](https://github.com/knative-sandbox/eventing-kafka)项目中, 由于错误的没有关闭，导致了一个http request无法有效复用连接的错误。在通过我的PR加入手动Close之后，QPS能够提升一倍以上：

https://github.com/knative-sandbox/eventing-kafka/pull/239

# 端口耗尽问题

当使用Go client不断请求主机连接，而又没有合适的设定`MaxConnsPerHost`时，TCP连接将会不断被创建，每创建一个TCP连接，就要用掉一个主机端口。当`/etc/sysctl.conf`中的`net.ipv4.ip_local_port_range`指定范围端口全部耗尽时，新建立连接将会报错:

```bash
connect: cannot assign requested address
```
解决该问题，应当妥善设置最大连接数限制，在应用程序中设置适当的并发数。

