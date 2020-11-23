---
title: "Go,Javascript与Websocket"
date: 2019-10-19T15:45:40+08:00
tags: "Go"
---


# JS中建立Websocket连接

```js
var ws = new WebSocket("ws://hostname/path", ["protocol1", "protocol2"])
```

## 参数说明

第一个参数是服务端websocket地址，如果是https+websocket，那么前缀写成wss

第二个参数并不是必须的，它约定了双方通讯使用的自定义子协议，会被放到这个Header中： Sec-WebSocket-Protocol

子协议在某些场合是很必要的，例如服务端要与多个客户端版本兼容，那么若干个版本之后，服务端设定支持子协议 v1.5, v2.0, 而客户端发送的却是 v1.0，那么他们就可以在握手阶段失败，不会继续通信下去导致奇奇怪怪的错误。

## 携带额外信息及认证

WebSocket构造函数只有两个变量，不能提供通过设置自定义Header的方式来携带其它信息，但仍可以通过一些取巧的办法携带额外的信息，用于认证等：

1. 通过ws地址填写形如 ws://username:password@hostname/path， 即构造出了 Authorization Header

2. 通过ws地址填写形如 ws://:password@hostname/path ，即构造出了 Bearer Token Header

3. 通过在Cookie中加入值，也能够携带额外的信息

因此，在服务端设计握手阶段认证时，应当避免使用这三种方式外携带的信息来进行认证（例如设置一个自定义的头部），当然也可以在websocket连接建立后，再通过自定义的认证协议，走websocket进行认证。

# Go中提供Websocket服务

Google自己提供一个Websocket包 : golang.org/x/net/websocket

不过他们亲口承认这个包缺乏一些特性，也缺乏维护，他们推荐用 github.com/gorilla/websocket (原文见 https://godoc.org/golang.org/x/net/websocket)

```go
// 这里代码使用了go-restful 作为http框架，换成http也无妨
conn, err := websocket.Upgrade(resp.ResponseWriter, req.Request, nil, 0, 0)
if err != nil {
    resp.WriteError(http.StatusBadRequest, err)
    return
}
defer conn.Close()
```

也可以手动创建一个 Upgrader 来处理子协议协商问题, 如果协商通过，就可以很容易的获得最终协商好的子协议，从而使用正确版本的数据格式和处理方法。

# 数据发送

WebSocket发送的数据都是“帧（Frames）”，主要有这么几种：

* 持续帧（用于数据分片，一般不明确使用）
* 文本帧（传输文本数据）
* 二进制帧（传输二进制数据）
* Ping/Pong帧 （用于心跳等，简单检测连接存活状态）
* 控制帧（关闭连接等）

JS中提供了send方法，能够发送文本帧或二进制帧：
```js
ws.send('{"abc":"def"}')
```
通过调用 ws.close(code, msg)， 可以发送关闭信息，如果不提供，那么默认code为1005（正常关闭），而不明确关闭，那么服务端收到的可能是1006.

当服务端对连接发起Ping时，浏览器中活跃的WebSocket对象会自动回复Pong, 这可以用于连接的活跃检测。

服务端向客户端发送的数据就要自由的多了，在此不多讲，参考包文档即可。

# 数据接收

```js
ws.onmessage = function(event) {
    graphData = JSON.parse(event.data)
    console.log('Received graph data:', graphData)
    if(graphData.error != null) {
        loadGraphData(null, graphData.error)
        return
    }
    loadGraphData(graphData, 'success')
};
```

在服务端的数据接收一般需要一个单独的go routine进行处理，可以使用 NextReader, ReadMessage, ReadJSON这几个方法进行读取。需要注意的是，对于同一个WebSocket连接， 这些读取方法应当在同一个go routine中顺序执行，否则读取操作将导致上一个进行中的读取失败。

# WebSocket的关闭

JS和Go中都提供了WebSocket的关闭事件监听，如:
```js
ws.onclose = function(event) {
    if (event.wasClean) {
        //alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
        console.log('WSError:', event)
        showError('服务暂不可用，请稍后刷新')
    }
};
```
```go
conn.SetCloseHandler(func(code int, msg string) error {
    c.closed <- msg
    return nil
})
```
但是需要知道的是，WebSocket的关闭事件是**基于收到明确的关闭消息**情况下才会出现的。换言之， 你在服务端监听WebSocket关闭事件，当浏览器页面刷新或关闭时，服务端并不能及时发现旧的WebSocket已经关闭，甚至向它们发送信息仍然不会收到任何错误（浏览器不关闭的情况下）。
同样的，当服务器突然关闭，而关闭前又没有明确关闭所有WebSocket连接时，那么在JS中写的onclose事件也不会起作用。

因此，服务端要想避免无用的WebSocket占用资源，应当维护一种心跳机制，而WebSocket协议已经提供了Ping/Pong帧用于做这件事，而JS中的WebSocket对象也默认能够回应Ping帧，因此心跳方案是：

* 在服务端建立WebSocket连接后，每隔一个心跳周期T，向连接发送Ping，设定回应时限小于T（建议设置为 T/2)
* 当Ping返回错误时，表明客户端已离线，或因种种原因未能在时限内回应Ping，服务端关闭连接，并将其移除

```go
heartbeat := time.Duration(s.cfg.Heartbeat) * time.Second
tick := time.NewTicker(heartbeat)
for {
    select {
    case <-tick.C:
        err := c.WriteControl(websocket.PingMessage, []byte(key), time.Now().Add(heartbeat/2))
        if err != nil {
            log.Printf("Websocket Idle connection %s: Ping received error %s", key, err.Error())
            return
        }
    case <-s.ctx.Done():
        log.Printf("Websocket closed for client %s: server is closing\n", key)
        return
    case msg := <-c.closed:
        log.Printf("Websocket closed for client %s: %s\n", key, msg)
        return
    }
}
```

这种方案的优点是利用了JS自带的Pong回应，不需要写额外的JS代码，而服务端实现也比较简单。

先前网上查了一些方案，是在客户端借助send发送数据帧来进行心跳，同样的服务端也要单独针对这种数据帧进行处理，有些过于复杂了。

# JS防止堵塞

WebSocket连接建立后，客户端与服务端之间建立了一条长连接，其后所有的数据通讯都要在这一个socket上传输。当客户端频繁发送数据时，数据就可能会堵塞在本地，JS中提供了一种方法可以限制发送速率：
```js
// 每隔100毫秒发送数据， 这里限制为只有在没有缓存时才会发送数据
setInterval(() => {
  if (socket.bufferedAmount == 0) {
    socket.send(moreData());
  }
}, 100);
```
通过读取 bufferedAmount 就可以知道当前连接是否堵塞， 可以决定是否继续发送。

# JS中的错误处理

握手阶段出现错误：

如果由于子协议协商等原因，WebSocket未能成功升级，那么会在 onclose 事件中收到消息

其它错误在ws.onerror 中处理

ws.onclose 事件接收到的 event主要有 wasClean(bool), code(int), reason(string) 几个成员，主要作用是：

wasClean 当js主动发起ws.close时，onclose接收到的事件中该值为true，其他情况不管服务器以什么状态码关闭，是false

code: 关闭时的状态码

reason: 关闭时发送的关闭消息，一般用于详细说明情况（如：服务器正在关闭/重启，等等）

# 关闭

在服务端调用 websocket.Conn.Close() 的作用是关闭底层socket， **不能**让client收到有效的消息，正确的关闭应当手动发送Close消息：
```go
conn.WriteControl(websocket.CloseMessage,
    websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server is closing"),
    time.Now().Add(time.Second))
```

同理，在JS客户端，不需要使用websocket时应当手动关闭WebSocket，页面刷新或离开时也应当关闭WebSocket：
```html
<body onunload="leavePage()">
...
<script>
    function leavePage() {
        ws.close(1000, 'Client closed')
    }
</script>
```

WebSocket关闭的常用状态码，参见 websocket/conn.go 或RFC文档，几个常见的WebSocket关闭码：

* 1000: 正常关闭
* 1001: 服务端暂停服务，或客户端离开页面
* 1005: 这次关闭没有状态码
* 1006: 这是一次不正常关闭，没有发送关闭帧

# 网关配置

在生产实践中, 会发现经常出现写入Websocket失败:

```
write tcp xxx.xxx.xxx.xxx:9000-> xxx.xxx.xxx.xxx:xxxx: write: broken pipe
```

broken pipe 错误见于向一个已经关闭的socket写入数据. 但是服务端已经通过心跳确保客户端连接存活, 客户端也没有主动切断连接, 为什么还会出现这种错误呢?

后来发现是网关nginx默认配置了连接超时, 超过这个时间, 连接就会被干掉.

解决方案:

1. 网关对websocket延长超时

上调 ```proxy_read_timeout```, 默认为60秒

2. 客户端增加重连机制

当客户端发现无法连接到服务端时, 主动重连服务端. 因为在这种时候, 服务端不可能主动反向重连.

3. 服务端有多实例, 而Websocket有状态时, 配置网关转发规则, 确保切断后仍能连回同一服务端.

# 参考文档

JS： https://javascript.info/websocket#opening-a-websocket

Go： https://godoc.org/github.com/gorilla/websocket

WebSocket： https://github.com/HJava/myBlog/tree/master/WebSocket%20%E5%8D%8F%E8%AE%AE%20RFC%20%E6%96%87%E6%A1%A3