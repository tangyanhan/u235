---
title: "DNS懒人包"
date: 2021-01-19T10:58:43+08:00
draft: true
---

一般的，DNS(Domain Name System)，是为了解决IP地址难记的问题而提出的解决方案，提供从域名到若干IP地址的转换。

现存的DNS当然不只有这些功能，本文尝试把现实生活中可能会遇到的DNS相关问题与知识点结合起来。

## DNS 服务器

DNS服务器存放了DNS域名系统的数据库，收到客户端请求时可以进行对应的查询，返回结果。
Linux下的DNS服务器设置存放于 `/etc/resolv.conf`中，这个文件是自动生成的，修改这个文件的效果，在重启后会丢失。
可以通过如`systemd-resolve --set-dns=114.114.114.114 --interface=enp0s31f6`这种命令增加dns，或者通过[这些方法](https://devilbox.readthedocs.io/en/latest/howto/dns/add-custom-dns-server-on-linux.html)暂时或永久的改变DNS服务器。


```bash
ethan@ethan:/etc$ cat resolv.conf 
nameserver 223.5.5.5
nameserver 8.8.8.8


nameserver 127.0.0.53
```

### 域名覆盖

DNS服务器提供了一个动态查询域名对应IP的方式，在我们自己制定域名时，也可以直接指定某个IP为某个域名，这时可以改写 `/etc/hosts`文件,在其中添加记录：

```
127.0.0.1       localhost
127.0.1.1       ethan
192.168.99.101  node-1
```

在Kubernetes中，对于Pod，可以通过指定`spec.hostAliases`增加记录：

```yaml
hostAliases:
- hostnames:
  - node-1
  - node-one
  ip: 192.168.99.101
```

但这需要每个需要解析的部署都增加配置才行。我们也可以修改DNS，这样所有的Pod只要通过CoreDNS进行域名解析，就可以起到域名覆盖的效果了。

如果使用CoreDNS，可以通过`hosts`[配置](https://coredns.io/plugins/hosts/)，起到类似hosts的效果：

```
        hosts /etc/myhosts {
           fallthrough
        }
```
然后挂载到CoreDNS中一个 `/etc/myhosts`文件即可。

## 缓存与查询顺序

DNS进行一次解析后，会在本地内存中记录缓存，在缓存过期之前，并不需要重新请求DNS，这样可以起到提高性能，减少查询的作用。

glibc进行DNS查询时，会通过`/etc/nsswitch.conf`控制查询的顺序：

```
hosts:          files mdns4_minimal [NOTFOUND=return] dns mdns4
```

[这一行的含义是:](https://ubuntu.com/server/docs/network-configuration)

* `files` 意思是首先找`/etc/hosts`中的静态域名

* `mdns4_minimal` 尝试通过多播DNS解析域名

* `[NOTFOUND=return]` 意味着如果前面的步骤都得不到应答，应该认为这个结果就是无解，不应该继续下去

* `dns`含义时通过旧的单播DNS查询域名

* `mdns4` 含义是多播DNS查询

然而有些特殊的容器镜像，如`alpine`镜像，使用musl libc代替了glibc，`alpine`镜像中默认是没有`nsswitch.conf`的。

Golang静态编译的程序不使用glibc链接库时，会使用自己的逻辑来进行DNS查询。首先它会查询`/etc/nsswitch.conf`，如果找不到，那么它会默认使用DNS，然后使用文件。即DNS优先，`/etc/hosts`往后排。

当一个Go静态编译的程序放入`alpine`镜像时，由于找不到`/etc/nsswitch.conf`，此时就会优先使用DNS查询，查询失败才会查询`/etc/hosts`。这样会导致每次总会有一定的DNS查询延迟。

在[这篇文章](https://mp.weixin.qq.com/s/VYBs8iqf0HsNg9WAxktzYQ)中描述了一个`conntrack`导致DNS查询丢包的问题，表面现象就是通过DNS查询时，偶发几秒钟的延迟。


## 查询

DNS查询是一种构建于传输层之上的协议，可以使用TCP或UDP传输。在使用glibc进行DNS查询时，通过`man resolf.conf`可知，默认具有5秒的超时时间。Golang目前总是使用TCP进行DNS查询(dnsclient_unix.go)。

glibc在进行域名查询时，会同时发起`A`(IPv4)和`AAAA`(IPv6)两种记录的解析请求，在`musl libc`的实现中，有些特定机器下会出现并发请求id冲突，而被拦截的情况，从而导致请求失败。[DNS解析异常](https://www.bookstack.cn/read/kubernetes-practice-guide/troubleshooting-cases-dns-resolution-abnormal.md)。这个问题在不使用`musl libc`时可以得到解决，方法是不使用`alpine`基础镜像，如go程序，可以禁用cgo，从而使用golang自己的逻辑，通过TCP并发请求。

## 资源记录类型

DNS通常用来确认域名对应的IP地址，但也可以用于从IP地址找到域名，从域名找到别名，获得IPv4和IPv6地址，以及其它非互联网数据提供分布式数据库功能。

使用`nslookup`命令进入交互模式，通过`set type=<>`可以控制DNS解析资源类型

* A/AAAA型记录

`A`型记录为IPv4地址，`AAAA`型记录为IPv6地址。在使用glibc进行DNS解析时，会并发发出这两种类型的查询。在阿里云的DNS解析中，也可以看到这类的记录

![Aliyun](/images/post/Screenshot_20210119_104944.png)

* CNAME记录

CNAME即规范名称(Canonical name)，会将单一域名的别名记录进系统中。很多网站为了避免用户typo进错网站，或者出于其它目的，会将一组域名关联到自己的主机上。

```bash
ethan@ethan:~$ host -t any cn.bing.com
cn.bing.com is an alias for cn-bing-com.cn.a-0001.a-msedge.net.
```

* 逆向DNS查询: PTR指针记录

PTR记录查询可以反向找出对应域名。由于域名是域名服务器树状结构构成的，需要一层层取解析。可以通过`.in-addr.arpa.`域来找出对应。


