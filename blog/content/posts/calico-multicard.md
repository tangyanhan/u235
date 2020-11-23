---
title: "Calico 多网卡解决方式"
date: 2020-03-30T16:19:15+08:00
tags: "Kubernetes"
---

当机器内存在多张网卡时, Calico的daemonset, 即 calico-node 可能使用错误的网卡尝试与其它主机建立连接.
以VirtualBox创建的多节点虚拟机为例, 我们一般需要两个网卡, 一个NAT网卡用于联系外网, 另一个虚拟网卡用于与host和其它节点建立一个虚拟局域网. 如果这时calico错误使用了NAT网卡, 那么自然无法连上其它主机.

通过调整calicao 网络插件的网卡发现机制，修改 IP_AUTODETECTION_METHOD 对应的value值。官方提供的yaml文件中，ip识别策略（IPDETECTMETHOD）没有配置，即默认为first-found，这可能会导致一个网络异常的ip作为nodeIP被注册，从而影响 node-to-node mesh 。我们可以修改成 can-reach 或者 interface 的策略，尝试连接某一个Ready的node的IP，以此选择出正确的IP。

这个环境变量需要添加到 calico.yaml 中, 位置在 ```DaemonSet/calico-node.spec.template.spec.containers[0].env```

can-reach使用您的本地路由来确定将使用哪个IP地址到达提供的目标。可以使用IP地址和域名。
```
# Using IP addresses
IP_AUTODETECTION_METHOD=can-reach=8.8.8.8
IP6_AUTODETECTION_METHOD=can-reach=2001:4860:4860::8888
```

```
# Using domain names
IP_AUTODETECTION_METHOD=can-reach=www.google.com
IP6_AUTODETECTION_METHOD=can-reach=www.google.com
```

interface使用提供的接口正则表达式（golang语法）枚举匹配的接口并返回第一个匹配接口上的第一个IP地址。列出接口和IP地址的顺序取决于系统。

```
# Valid IP address on interface eth0, eth1, eth2 etc.
IP_AUTODETECTION_METHOD=interface=eth.*
IP6_AUTODETECTION_METHOD=interface=eth.*
```

修改后, 重新apply即可.
