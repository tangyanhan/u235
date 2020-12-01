---
title: "Linux组网基础:虚拟设备"
date: 2020-12-01T23:20:34+08:00
---

# 基本概念

在讨论Linux网络设备时，一般会涉及“层”的概念，例如L1,L2,L3等等。这时使用的是OSI七层模型，即“物数网传会表应”。

网桥(bridge)：工作在L2（数据链路层），只有目的MAC匹配的数据包才能被发送到出口，只能在同一网段下进行流量转发。

交换机(switch): 工作在L2（也可能包含L3），可以将多组MAC根据不同的规则进行转发。

路由(route): 工作在L3(网络层)，是基于IP规则转发的，因而能够连接不同子网，也是当前互联网能够互联的基础。

# 网络命名空间:Network Namespace

相关命令: `ip netns`

相关链接: https://www.man7.org/linux/man-pages/man8/ip-netns.8.html

## 什么是网络命名空间

网络命名空间是Linux内核网络桟的一个副本，其中包含了自己的路由(route)，防火墙规则(iptable)和网络设备(link)。


默认情况下，进程从父进程继承网络空间，因此一开始所有进程都从init进程中继承同样的网络空间。

一个有名字的网络命名空间，在文件系统中则是`/var/run/netns/NAME`，是可以打开的。 通过打开它获得的fd就指向了那个网络命名空间。持有fd也就能够让命名空间保持存活。 这个fd可以通过[`setns`](https://www.man7.org/linux/man-pages/man2/setns.2.html)系统调用来改变一个task的网络命名空间。

## 基本操作

`ip netns list` 列出网络命名空间

`ip netns add NAME` 创建一个新的网络命名空间

`ip netns delete NAME` 删除网络命名空间

`ip netns pids NAME` 列出网络命名空间下的所有进程pid

`ip netns identify [PID]` 找出进程所在的网络空间名称

`ip netns monitor` 将会对接下来的netns操作进行监控，添加和删除操作会被记录在日志里

`ip netns exec`会在指定网络命名空间中运行一个新的进程。

`ip netns attach NAME PID` 将指定的PID挂载到指定的ns中，如果ns不存在，那么创建一个。

# 路由(Route)

相关链接：

[The IPv4 Routing Subsystem](https://apprize.best/linux/kernel/6.html)

[What is IP routing](https://study-ccna.com/what-is-ip-routing/)

路由是连接不同子网的网络设备，可以是硬件的，也可以是软件的。它通过不同的规则，将来自一个子网的信息转发到另一个子网中。

任意两个IP地址相互访问时，如何从一个到达另一个就成了问题。路由通过路由协议（如OSPF,RIP，BGP）等维护路由表，这些路由表就是一个个路口道标一样，通过它们不断的跳转转发，经过若干跳(hop)之后，到达目的地。

## 默认网关(Default Gateway)

主机A现在有路由器R1的IP地址作为默认网关，当A想要访问一个远程地址B，又不知道怎么去连接它，那么A就可以将数据包发给默认网关R1，期待着R1拥有访问B的途径。

![DefaultGateway](https://x9s4w2e2.stackpathcdn.com/wp-content/uploads/2016/01/default_gateway.jpg)

## 路由表(Routing Table)

每个路由器，包括每台主机自身，都会在RAM中维护一张路由表，主要包含以下内容：

* 目的地网络及子网掩码：定义了一个范围的网络
* 远程路由器(remote router): 访问该网络的路由，要发向该网络，那么就将数据包转发到这里
* 出口网卡(outgoing interface): 要将流量转发到该网络，流量出口应该是哪个网卡


可以通过`ip route show`命令列出当前默认网络命名空间下的route：

```
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 
```

`default`路由在其它路由匹配不到时，就会将流量发送到default。

`ip route add default via 192.168.99.112`

```
ethan@ethan:~$ ip route show
default via 192.168.30.254 dev enp0s31f6 proto dhcp metric 100 
default via 192.168.8.1 dev wlp4s0 proto dhcp metric 600 
169.254.0.0/16 dev docker0 scope link metric 1000 linkdown 
172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1 linkdown 
172.18.0.0/16 dev br-9cb479a86720 proto kernel scope link src 172.18.0.1 linkdown 
192.168.8.0/24 dev wlp4s0 proto kernel scope link src 192.168.8.108 metric 600 
192.168.30.0/24 dev enp0s31f6 proto kernel scope link src 192.168.30.197 metric 100 
192.168.99.0/24 dev vboxnet1 proto kernel scope link src 192.168.99.1 
```

一个纯粹的路由器工作在L3，通过路由表获得下一跳的IP地址，这个地址是可以通过MAC地址直连的，然后通过数据链路层将数据包转发到该地址即可。

但是我们实际场景中，并不只是公网IP不断的子网划分，或者是局域网的流量跳转到更小的子网中。例如我们请求访问baidu.com，我们的信息通过路由若干跳之后，成功抵达外网主机。写回数据包在软件层面上，通过socket写回，仿佛在向一条电缆简单的发数据，然后就抵达了请求者，但实际baidu.com的主机不可能通过外网的路由器知道我们家里的电脑源地址192.168.1.2这样的IP该怎么转发-现实中这样的子网成千上万。

类似的，主机上运行容器，每个容器都有一个容器IP。当我们从容器内访问外网时，外网流量写回网络接口时，如何决定将流量投递到那个请求的容器呢？

之所以路由能够解决这个问题，是因为我们实际使用的路由并非纯粹的路由器。通过地址变换，把流量在两个完全不同的网络中转发，这就是NAT(Network Address Transform). 在公有云中，我们常见`DNAT`(Desitionation NAT)和`SNAT`(Source NAT)，DNAT转换目的地址使流量能够从外面进入子网从而提供服务，SNAT转换源地址，使得内部流量能够流出访问外网。我们自家的网关默认都是具备SNAT,而公有云则一般分开管理。

# 网桥(Network Bridge)

网桥是一个工作在L2的交换机，提供了若干个网线接口，可以把不同的设备连接起来，使之在L2成为一个互联的网络。可以是虚拟的，也可以是实体的。

##　通过ip命令操作

添加一个网桥，名为br1，并将其启动:

```
sudo ip link add dev br1 type bridge
sudo ip link set dev br1 up
```

这时，系统中已经存在了一个文件夹`br1`:
```
/sys/devices/virtual/net/br1
```

通过`ip link show br1`，会发现其状态从`DOWN`变成了`UNKNOWN`.

这是因为网桥还没有绑定状态为`UP`的网卡。将有线接口设定为网桥br1的子接口：

```
sudo ip link set dev enp0s31f6 master br1
```

这时再次查看br1状态，会发现状态已经变成了UP。这时，如果通过以下命令将网卡关闭，就会看到网桥变成DOWN：
```
sudo ip link set dev enp0s31f6 down
```
将有线网卡重新起来，就会发现网桥状态也重新UP。

我们的有线接口`enp0s31f6`插入到了网桥br1中，这时任何连接到网桥br1的设备同样可以获得来自`enp0s31f6`的流量了。但是，这个前提是我们的有线网卡必须打开混杂模式(PROMISC)，否则，网卡会自动过滤掉跟自己MAC地址不同的数据包:

```
sudo ip link set enp0s31f6 promisc on
```

给网桥添加IP地址:
```
ip addr add dev bridge_name 192.168.66.66/24
```

删除网桥之前，务必先对端口进行还原：
```
# 关闭混杂模式
sudo ip link set enp0s31f6 promisc off
# 关闭端口
sudo ip link set enp0s31f6 down
# 从网桥移除端口
sudo ip link set dev enp0s31f6 nomaster
```

然后，可以删除网桥：
```
sudo ip link delete br1 type bridge
```

## 通过brctl操作

网桥操作的`brctl`命令需要提前安装软件包:
```
sudo apt install -y bridge-utils
```

添加一个网桥:
```
sudo brctl addbr ethan
```

添加完网桥，通过`ip addr`你会发现同时多出了一个同名的网卡。

列出目前所有的网桥:
```
brctl show
```

向网桥`ethan`添加网络接口`vboxnet0`:
```
sudo brctl addif ethan vboxnet0
```

使用网桥前需要先启动它，启动网桥:
```
sudo ip link set up dev ethan
```

删除网桥。你无法删除一个启动的网桥，需要先用ip link set dev NAME down来关闭它，然后才能删除:
```
ethan@ethan:~$ sudo brctl delbr ethan
bridge ethan is still up; can't delete it
ethan@ethan:~$ sudo ip link set dev ethan down
ethan@ethan:~$ sudo brctl delbr ethan
```

